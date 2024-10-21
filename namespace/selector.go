package namespace

import "github.com/shamcode/simd/storage"

type selectorOperation uint8

const (
	union selectorOperation = iota + 1
	intersection
)

type result struct {
	items               []storage.IDIterator
	idsUnique           bool
	size                int
	operation           selectorOperation
	operationRecognized bool
}

type resultByBracketLevel map[int]*result

func (byLevel resultByBracketLevel) save(
	level int,
	items []storage.IDIterator,
	idsUnique bool,
	size int,
	operation selectorOperation,
) {
	res, exist := byLevel[level]
	if !exist {
		byLevel[level] = &result{
			items:               items,
			idsUnique:           idsUnique,
			size:                size,
			operation:           operation,
			operationRecognized: false,
		}
		return
	}

	switch operation {
	case union:
		res.size += size
		res.items = append(res.items, items...)
		res.idsUnique = false // TODO: optimize for id < 2 OR id > 5
	case intersection:
		if res.size > size {
			res.items = items
			res.size = size
			res.idsUnique = idsUnique
		}
	}
}

func (byLevel resultByBracketLevel) reduce(fromLevel int, toLevel int) ([]storage.IDIterator, int, bool, bool) {
	var items []storage.IDIterator
	var size int
	var item *result
	var idsUnique bool
	isFirst := true
	for fromLevel > toLevel {
		item = byLevel[fromLevel]
		if nil != item {
			if isFirst {
				items = item.items
				size = item.size
				idsUnique = item.idsUnique
				isFirst = false
			} else {
				switch item.operation {
				case union:
					size += item.size
					items = append(items, item.items...)
					idsUnique = false
				case intersection:
					if item.size < size {
						items = item.items
						size = item.size
						idsUnique = item.idsUnique
					}
				}
			}
			delete(byLevel, fromLevel)
		}
		fromLevel -= 1
	}
	return items, size, idsUnique, !isFirst
}
