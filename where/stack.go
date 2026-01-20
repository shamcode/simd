package where

import "sync"

// resultsByBracketLevel usage for merge results by condition's brackets level.
type resultsByBracketLevel map[int]*result

type result struct {
	value        bool
	isAnd        bool
	opRecognized bool
}

func getFromPool() *result {
	return bracketLevelResultPool.Get().(*result)
}

func returnToPool(a *result) {
	a.opRecognized = false
	bracketLevelResultPool.Put(a)
}

var bracketLevelResultPool = &sync.Pool{
	New: func() any {
		return &result{} //nolint:exhaustruct
	},
}

func (stack resultsByBracketLevel) save(level int, value, isAnd bool) {
	item, exists := stack[level]
	switch {
	case !exists:
		// A
		// First item, just save
		item = getFromPool()
		item.value = value
		stack[level] = item
	case isAnd:
		// ... AND A
		// Not first item, merge as AND
		item.value = item.value && value
	default:
		// ... OR A
		// Not first item, merge as OR
		item.value = item.value || value
	}
}

func (stack resultsByBracketLevel) reduce(fromLevel int, toLevel int) bool {
	stackItem := stack[fromLevel]
	result := stackItem.value
	returnToPool(stackItem)
	delete(stack, fromLevel)

	fromLevel -= 1
	for fromLevel > toLevel {
		stackItem := stack[fromLevel]
		if nil != stackItem {
			if stackItem.isAnd {
				// (...) AND A
				result = result && stackItem.value
			} else {
				// (...) OR A
				result = result || stackItem.value
			}

			returnToPool(stackItem)
			delete(stack, fromLevel)
		}

		fromLevel -= 1
	}

	return result
}

func (stack resultsByBracketLevel) pop(fromLevel int, toLevel int) {
	for fromLevel > toLevel {
		stackItem := stack[fromLevel]
		if nil != stackItem {
			returnToPool(stackItem)
			delete(stack, fromLevel)
		}

		fromLevel -= 1
	}
}
