package types

// Maybe represents "found T" vs "legitimately not present". Use it for lookups
// where absence is a valid outcome, not a failure. Errors still mean the lookup itself failed.
type Maybe[T any] struct {
	Value     T
	IsPresent bool
}

func Some[T any](v T) Maybe[T] { return Maybe[T]{Value: v, IsPresent: true} }

func None[T any]() Maybe[T] {
	var zero T
	return Maybe[T]{Value: zero, IsPresent: false}
}
