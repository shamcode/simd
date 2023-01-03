package indexes

import (
	"fmt"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/bytype"
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"log"
)

var _ executor.Namespace = (*NamespaceWithIndexes)(nil)

type NamespaceWithIndexes struct {
	logger  Logger
	storage *storage.RecordsByID
	byField map[string]*bytype.Index
}

func (ns *NamespaceWithIndexes) Get(id int64) record.Record {
	return ns.storage.Get(id)
}

func (ns *NamespaceWithIndexes) Insert(item record.Record) error {
	if nil != ns.Get(item.GetID()) {
		return fmt.Errorf("%w: ID == %d", ErrRecordExists, item.GetID())
	}
	ns.insert(item)
	return nil
}

func (ns *NamespaceWithIndexes) insert(item record.Record) {
	item.ComputeFields()
	id := item.GetID()
	for _, idx := range ns.byField {
		key := idx.Compute.ForRecord(item)
		records := idx.Storage.Get(key)
		if nil == records {
			records = storage.NewIDStorage()
			idx.Storage.Set(key, records)
		}
		records.Add(id)
	}
	ns.storage.Set(id, item)
}

func (ns *NamespaceWithIndexes) Delete(id int64) error {
	item := ns.Get(id)
	if nil == item {
		return nil
	}
	for _, idx := range ns.byField {
		records := idx.Storage.Get(idx.Compute.ForRecord(item))
		if nil != records {
			records.Delete(id)
		}
	}
	ns.storage.Delete(id)
	return nil
}

func (ns *NamespaceWithIndexes) Upsert(item record.Record) error {
	id := item.GetID()
	oldItem := ns.Get(id)

	if nil == oldItem {

		// It's insert
		ns.insert(item)
		return nil
	}

	// It's update
	item.ComputeFields()
	for _, idx := range ns.byField {

		oldValue := idx.Compute.ForRecord(oldItem)
		newValue := idx.Compute.ForRecord(item)

		if newValue == oldValue {

			// Field index not changed, ignore
			continue
		}

		// Remove old item from index
		oldRecords := idx.Storage.Get(oldValue)
		if nil != oldRecords {
			oldRecords.Delete(id)
		}

		records := idx.Storage.Get(newValue)
		if nil == records {

			// It's first item in index, create index storage
			records = storage.NewIDStorage()
			idx.Storage.Set(newValue, records)
		}

		// Add new item to index
		records.Add(id)
	}
	ns.storage.Set(id, item)
	return nil
}

func (ns *NamespaceWithIndexes) AddIndex(index *bytype.Index) {
	ns.byField[index.Field] = index
}

func (ns *NamespaceWithIndexes) SelectForExecutor(conditions where.Conditions) ([]record.Record, error) {
	byLevel := make(resultByBracketLevel)
	lastBracketLevel := 0

	for _, condition := range conditions {
		var op selectorOperation
		if condition.IsOr {
			op = union
		} else {
			op = intersection
		}

		if lastBracketLevel > 0 {
			last := byLevel[lastBracketLevel]
			if nil != last && condition.BracketLevel >= lastBracketLevel && !last.operationRecognized {
				last.operationRecognized = true
				last.operation = op
			}
		}

		exists, indexSize, ids, err := ns.getIndexForCondition(condition)
		if nil != err {
			return nil, err
		}
		if !exists {
			all := ns.storage.GetIDStorage()
			ids = append(ids, all)
			indexSize = ns.storage.Count()
		}

		if lastBracketLevel > condition.BracketLevel {
			subLevelItems, subLevelSize, hasItems := byLevel.reduce(lastBracketLevel, condition.BracketLevel)
			if hasItems {
				switch op {
				case union:
					ids = append(ids, subLevelItems...)
					indexSize += subLevelSize
				case intersection:
					if subLevelSize < indexSize {
						ids = subLevelItems
						indexSize = subLevelSize
					}
				}
			}
		}

		byLevel.save(condition.BracketLevel, ids, indexSize, op)
		lastBracketLevel = condition.BracketLevel
	}

	items, size, hasItems := byLevel.reduce(lastBracketLevel, 0)

	if !hasItems {
		ns.logger.Println("index not applied", conditions)
		return ns.storage.GetAllData(), nil
	}

	if size >= ns.storage.Count() {
		ns.logger.Println("index not applied (large select)", conditions)
		return ns.storage.GetAllData(), nil
	}
	return ns.storage.GetData(items), nil
}

func (ns *NamespaceWithIndexes) SetLogger(logger Logger) {
	ns.logger = logger
}

func (ns *NamespaceWithIndexes) getIndexForCondition(condition where.Condition) (
	indexExists bool,
	count int,
	ids []storage.LockableIDStorage,
	err error,
) {
	index, exists := ns.byField[condition.Cmp.GetField()]
	if !exists {
		return
	}
	indexExists = true

	cmpType := condition.Cmp.GetType()
	if cmpType == where.EQ && !condition.WithNot {

		// A == '1'
		itemsByValue := index.Storage.Get(index.Compute.ForComparatorFirstValue(condition.Cmp))
		if nil != itemsByValue {
			count = itemsByValue.Count()
			ids = append(ids, itemsByValue)
		}
		return
	}

	if cmpType == where.InArray && !condition.WithNot {

		// A IN ('1', '2', '3')
		index.Compute.EachComparatorValues(condition.Cmp, func(conditionValue interface{}) {
			itemsByValue := index.Storage.Get(conditionValue)
			if nil != itemsByValue {
				countForValue := itemsByValue.Count()
				if countForValue > 0 {
					count += countForValue
					ids = append(ids, itemsByValue)
				}
			}
		})
		return
	}

	keys := index.Storage.Keys()
	for _, key := range keys {
		resultForValue, errorForValue := index.Compute.Compare(key, condition.Cmp)
		if nil != errorForValue {
			err = errorForValue
			return
		}
		if condition.WithNot != resultForValue {
			count += index.Storage.Count(key)
			ids = append(ids, index.Storage.Get(key))
		}
	}

	return
}

func CreateNamespace() *NamespaceWithIndexes {
	return &NamespaceWithIndexes{
		logger:  log.Default(),
		storage: storage.NewRecordsByID(),
		byField: make(map[string]*bytype.Index),
	}
}
