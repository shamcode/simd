package namespace

import (
	"log"

	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type Namespace interface {
	Get(id int64) record.Record
	Insert(item record.Record) error
	Delete(id int64) error
	Upsert(item record.Record) error
	executor.Selector
}

type fieldsComputer interface {
	// ComputeFields is a special hook for optimize slow computing fields.
	// ComputeFields call on insert or update record.
	ComputeFields()
}

var _ Namespace = (*WithIndexes)(nil)

type WithIndexes struct {
	logger  Logger
	storage storage.RecordsByID
	indexes indexes.ByField
}

func (ns *WithIndexes) Get(id int64) record.Record {
	return ns.storage.Get(id)
}

func (ns *WithIndexes) Insert(item record.Record) error {
	if nil != ns.Get(item.GetID()) {
		return NewRecordAlreadyExists(item.GetID())
	}
	ns.insert(item)
	return nil
}

func (ns *WithIndexes) insert(item record.Record) {
	if item, ok := item.(fieldsComputer); ok {
		item.ComputeFields()
	}
	ns.indexes.Insert(item)
	ns.storage.Set(item.GetID(), item)
}

func (ns *WithIndexes) Delete(id int64) error {
	item := ns.Get(id)
	if nil == item {
		return nil
	}
	ns.indexes.Delete(item)
	ns.storage.Delete(id)
	return nil
}

func (ns *WithIndexes) Upsert(item record.Record) error {
	id := item.GetID()
	oldItem := ns.Get(id)

	if nil == oldItem {
		// It's insert
		ns.insert(item)
		return nil
	}

	// It's update
	if item, ok := item.(fieldsComputer); ok {
		item.ComputeFields()
	}
	ns.indexes.Update(oldItem, item)
	ns.storage.Set(id, item)
	return nil
}

func (ns *WithIndexes) AddIndex(index indexes.Index) {
	ns.indexes.Add(index)
}

func (ns *WithIndexes) PreselectForExecutor(conditions where.Conditions) ( //nolint:funlen,gocognit,cyclop
	[]record.Record,
	error,
) {
	byLevel := make(resultByBracketLevel)
	lastBracketLevel := 0

	for _, condition := range conditions {
		var operation selectorOperation
		if condition.IsOr {
			operation = union
		} else {
			operation = intersection
		}

		if lastBracketLevel > 0 {
			last := byLevel[lastBracketLevel]
			if nil != last && condition.BracketLevel >= lastBracketLevel && !last.operationRecognized {
				last.operationRecognized = true
				last.operation = operation
			}
		}

		indexExists, indexSize, ids, idsUnique, err := ns.indexes.SelectForCondition(condition)
		if nil != err {
			return nil, err
		}
		if !indexExists {
			all := ns.storage.GetIDStorage()
			ids = append(ids, all)
			indexSize = ns.storage.Count()
		}

		if lastBracketLevel > condition.BracketLevel {
			subLevelItems, subLevelSize, subLevelIDSUnique, hasItems := byLevel.reduce(lastBracketLevel, condition.BracketLevel)
			if hasItems {
				switch operation {
				case union:
					ids = append(ids, subLevelItems...)
					indexSize += subLevelSize
					idsUnique = false
				case intersection:
					if subLevelSize < indexSize {
						ids = subLevelItems
						indexSize = subLevelSize
						idsUnique = subLevelIDSUnique
					}
				}
			}
		}

		byLevel.save(condition.BracketLevel, ids, idsUnique, indexSize, operation)
		lastBracketLevel = condition.BracketLevel
	}

	items, size, idsUnique, hasItems := byLevel.reduce(lastBracketLevel, 0)

	if !hasItems {
		ns.logger.Println("index not applied", conditions)
		return ns.storage.GetAllData(), nil
	}

	if size >= ns.storage.Count() {
		ns.logger.Println("index not applied (large select)", conditions)
		return ns.storage.GetAllData(), nil
	}
	return ns.storage.GetData(items, size, idsUnique), nil
}

func (ns *WithIndexes) SetLogger(logger Logger) {
	ns.logger = logger
}

func CreateNamespace() *WithIndexes {
	return &WithIndexes{
		logger:  log.Default(),
		storage: storage.CreateRecordsByID(),
		indexes: indexes.CreateByField(),
	}
}
