package namespace

import "github.com/shamcode/simd/storage"

type selectorOperation uint8

const (
	union selectorOperation = iota + 1
	intersection
)

type result struct {
	items               []storage.LockableIDStorage
	size                int
	operation           selectorOperation
	operationRecognized bool
}

type resultByBracketLevel map[int]*result

func (byLevel resultByBracketLevel) save(level int, items []storage.LockableIDStorage, size int, op selectorOperation) {
	res, exist := byLevel[level]
	if !exist {
		byLevel[level] = &result{
			size:      size,
			items:     items,
			operation: op,
		}
		return
	}

	switch op {
	case union:
		res.size += size
		res.items = append(res.items, items...)
	case intersection:
		if res.size > size {
			res.items = items
			res.size = size
		}
	}
}

func (byLevel resultByBracketLevel) reduce(fromLevel int, toLevel int) ([]storage.LockableIDStorage, int, bool) {
	var items []storage.LockableIDStorage
	var size int
	var item *result
	isFirst := true
	for fromLevel > toLevel {
		item = byLevel[fromLevel]
		if nil != item {
			if isFirst {
				items = item.items
				size = item.size
				isFirst = false
			} else {
				switch item.operation {
				case union:
					size += item.size
					items = append(items, item.items...)
				case intersection:
					if item.size < size {
						items = item.items
						size = item.size
					}
				}
			}
			delete(byLevel, fromLevel)
		}
		fromLevel -= 1
	}
	return items, size, !isFirst
}
