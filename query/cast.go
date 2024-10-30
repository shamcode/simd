package query

// Cast casts a value to a different type.
// If cast fails, an CastError is returned.
func Cast[From, To any](v From) (To, error) {
	res, ok := any(v).(To)
	if !ok {
		return res, CastError[To, From]{Expected: res, Actual: v}
	}

	return res, nil
}
