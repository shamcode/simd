package namespace

import (
	"log"

	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type Namespace[R record.Record] interface {
	Get(id int64) (R, bool)
	Insert(item R) error
	Delete(id int64) error
	Upsert(item R) error
	executor.Selector[R]
}

type fieldsComputer interface {
	// ComputeFields is a special hook for optimize slow computing fields.
	// ComputeFields call on insert or update record.
	ComputeFields()
}

type WithIndexes[R record.Record] struct {
	logger  Logger
	storage storage.RecordsByID[R]
	indexes indexes.ByField[R]
}

func (ns *WithIndexes[R]) Get(id int64) (R, bool) {
	return ns.storage.Get(id)
}

func (ns *WithIndexes[R]) Insert(item R) error {
	if _, exists := ns.Get(item.GetID()); exists {
		return NewRecordAlreadyExists(item.GetID())
	}

	ns.insert(item)

	return nil
}

func (ns *WithIndexes[R]) insert(item R) {
	if item, ok := any(item).(fieldsComputer); ok {
		item.ComputeFields()
	}

	ns.indexes.Insert(item)
	ns.storage.Set(item.GetID(), item)
}

func (ns *WithIndexes[R]) Delete(id int64) error {
	item, exists := ns.Get(id)
	if !exists {
		return nil
	}

	ns.indexes.Delete(item)
	ns.storage.Delete(id)

	return nil
}

func (ns *WithIndexes[R]) Upsert(item R) error {
	id := item.GetID()
	oldItem, exists := ns.Get(id)

	if !exists {
		// It's insert
		ns.insert(item)
		return nil
	}

	// It's update
	if item, ok := any(item).(fieldsComputer); ok {
		item.ComputeFields()
	}

	ns.indexes.Update(oldItem, item)
	ns.storage.Set(id, item)

	return nil
}

func (ns *WithIndexes[R]) AddIndex(index indexes.Index[R]) {
	ns.indexes.Add(index)
}

func (ns *WithIndexes[R]) PreselectForExecutor(conditions where.Conditions[R]) ( //nolint:funlen,cyclop
	[]R,
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

func (ns *WithIndexes[R]) SetLogger(logger Logger) {
	ns.logger = logger
}

func CreateNamespace[R record.Record]() *WithIndexes[R] {
	return &WithIndexes[R]{
		logger:  log.Default(),
		storage: storage.CreateRecordsByID[R](),
		indexes: indexes.CreateByField[R](),
	}
}
