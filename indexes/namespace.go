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

		indexExists, indexSize, ids, err := ns.selectFromIndexForCondition(condition)
		if nil != err {
			return nil, err
		}
		if !indexExists {
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

func (ns *NamespaceWithIndexes) selectFromIndexForCondition(condition where.Condition) (
	indexExists bool,
	count int,
	ids []storage.LockableIDStorage,
	err error,
) {
	var index *bytype.Index
	index, indexExists = ns.byField[condition.Cmp.GetField()]
	if !indexExists {
		return
	}
	if !condition.WithNot {
		switch condition.Cmp.GetType() {
		case where.EQ: // A == '1'
			count, ids = ns.selectFromIndexForEqual(index, condition)
			return
		case where.InArray: // A IN ('1', '2', '3')
			count, ids = ns.selectFromIndexForInArray(index, condition)
			return
		}
	}
	count, ids, err = ns.selectFromIndexForOther(index, condition)
	return
}

func (ns *NamespaceWithIndexes) selectFromIndexForEqual(index *bytype.Index, condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	itemsByValue := index.Storage.Get(index.Compute.ForValue(condition.Cmp.ValueAt(0)))
	if nil != itemsByValue {
		count = itemsByValue.Count()
		ids = append(ids, itemsByValue)
	}
	return
}

func (ns *NamespaceWithIndexes) selectFromIndexForInArray(index *bytype.Index, condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	for i := 0; i < condition.Cmp.ValuesCount(); i++ {
		itemsByValue := index.Storage.Get(condition.Cmp.ValueAt(i))
		if nil != itemsByValue {
			countForValue := itemsByValue.Count()
			if countForValue > 0 {
				count += countForValue
				ids = append(ids, itemsByValue)
			}
		}
	}
	return
}

func (ns *NamespaceWithIndexes) selectFromIndexForOther(index *bytype.Index, condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	keys := index.Storage.Keys()
	for _, key := range keys {
		resultForValue, errorForValue := index.Compute.Check(key, condition.Cmp)
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
