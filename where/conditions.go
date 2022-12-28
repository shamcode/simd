package where

import (
	"fmt"
	"github.com/shamcode/simd/record"
)

type Condition struct {
	WithNot      bool
	IsOr         bool
	BracketLevel int
	Cmp          FieldComparator
}

func (c Condition) String() string {
	return fmt.Sprintf("{%t %t %d %s %d}", c.WithNot, c.IsOr, c.BracketLevel, c.Cmp.GetField(), c.Cmp.GetType())
}

type Conditions []Condition

// Check checks that the record satisfies all the conditions
func (w Conditions) Check(item record.Record) (bool, error) {
	stack := make(resultsByBracketLevel)
	lastBracketLevel := 0

	for _, condition := range w {
		isAnd := !condition.IsOr

		if lastBracketLevel > 0 {
			last := stack[lastBracketLevel]
			if condition.BracketLevel >= lastBracketLevel {
				if !last.opRecognized {
					// A AND B
					// A OR B
					// A AND (B ... )
					// A OR (B ... )
					last.isAnd = isAnd
					last.opRecognized = true
				}
				if last.isAnd != last.value && last.isAnd == isAnd {
					// lazy or/and
					// skip Condition
					continue
				}
			}
		}

		compareResultForItem, err := condition.Cmp.Compare(item)
		if nil != err {
			return false, err
		}

		// Expression A != B it's analog for (A && !B) || (!A && B)
		conditionResult := condition.WithNot != compareResultForItem

		if lastBracketLevel > condition.BracketLevel {
			// ( ... ) AND B
			// ( ... ) OR B
			if isAnd != conditionResult {
				// Lazy reduce:
				// ( ... ) AND False == False
				// ( ... ) OR True == True
				stack.pop(lastBracketLevel, condition.BracketLevel)
			} else {
				subBracketsResult := stack.reduce(lastBracketLevel, condition.BracketLevel)
				if isAnd {
					conditionResult = subBracketsResult && conditionResult
				} else {
					conditionResult = subBracketsResult || conditionResult
				}
			}
		}

		stack.save(condition.BracketLevel, conditionResult, isAnd)

		lastBracketLevel = condition.BracketLevel
	}

	if 0 == len(stack) {
		return true, nil
	}

	return stack.reduce(lastBracketLevel, 0), nil
}
