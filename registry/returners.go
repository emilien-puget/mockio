package registry

import (
	"github.com/ovechkin-dm/mockio/matchers"
)

func ToReturner1[T any](retAll matchers.ReturnerAll) matchers.Returner1[T] {
	return &returner1impl[T]{
		all: retAll,
	}
}

func ToReturnerE[T any](retAll matchers.ReturnerAll) matchers.ReturnerE[T] {
	return &returnerEImpl[T]{
		all: retAll,
	}
}

type returnerAllImpl struct {
	methodMatch *methodMatch
	ctx         *mockContext
}

type returner1impl[T any] struct {
	all matchers.ReturnerAll
}

func (r *returner1impl[T]) ThenReturn(value T) matchers.Returner1[T] {
	return r.ThenAnswer(func(args []any) T {
		return value
	})
}

func (r *returner1impl[T]) ThenAnswer(f func(args []any) T) matchers.Returner1[T] {
	all := r.all.ThenAnswer(func(args []any) []any {
		return []any{f(args)}
	})
	return &returner1impl[T]{
		all: all,
	}
}

type returnerEImpl[T any] struct {
	all matchers.ReturnerAll
}

func (r *returnerEImpl[T]) ThenReturn(value T, err error) matchers.ReturnerE[T] {
	return r.ThenAnswer(func(args []any) (T, error) {
		return value, err
	})
}

func (r *returnerEImpl[T]) ThenAnswer(f func(args []any) (T, error)) matchers.ReturnerE[T] {
	all := r.all.ThenAnswer(func(args []any) []any {
		t, e := f(args)
		return []any{t, e}
	})
	return &returnerEImpl[T]{
		all: all,
	}
}

func (r *returnerAllImpl) ThenReturn(values ...any) matchers.ReturnerAll {
	return r.ThenAnswer(makeReturnFunc(values))
}

func (r *returnerAllImpl) ThenAnswer(f matchers.Answer) matchers.ReturnerAll {
	wrapper := &answerWrapper{
		ans: f,
	}
	r.methodMatch.addAnswer(wrapper)
	return r
}

func makeReturnFunc(values []any) matchers.Answer {
	return func(args []any) []interface{} {
		return values
	}
}

func NewReturnerAll(ctx *mockContext, data *methodMatch) matchers.ReturnerAll {
	return &returnerAllImpl{
		methodMatch: data,
		ctx:         ctx,
	}
}