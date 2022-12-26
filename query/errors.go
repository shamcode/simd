package query

import "errors"

var (
	ErrOrBeforeAnyConditions   = errors.New(".Or() before any condition not supported, add any condition before .Or()")
	ErrNotOpenBracket          = errors.New(".Not().OpenBracket() not supported")
	ErrCloseBracketWithoutOpen = errors.New("close bracket without open")
	ErrInvalidBracketBalance   = errors.New("invalid bracket balance: has not closed bracket")
)
