package useful

import (
	. "github.com/SimonRichardson/wishful/wishful"
)

type Validation interface {
	Of(v Any) Point
	Ap(v Applicative) Applicative
	Chain(f func(v Any) Monad) Monad
	Concat(y Semigroup) Semigroup
	Map(f Morphism) Functor
	Fold(f Morphism, g Morphism) Any
	Bimap(f Morphism, g Morphism) Monad
}

type failure struct {
	x Any
}

type success struct {
	x Any
}

func NewFailure(x Any) failure {
	return failure{
		x: x,
	}
}

func NewSuccess(x Any) success {
	return success{
		x: x,
	}
}

func (x failure) Of(v Any) Point {
	return NewSuccess(v)
}

func (x success) Of(v Any) Point {
	return NewSuccess(v)
}

func (x failure) Ap(v Applicative) Applicative {
	return v.(Validation).Fold(
		func(y Any) Any {
			return NewFailure(concatAnyvals(x.x)(y))
		},
		func(y Any) Any {
			return NewFailure(x.x)
		},
	).(Applicative)
}

func (x success) Ap(v Applicative) Applicative {
	return v.(Functor).Map(func(g Any) Any {
		fun := NewFunction(x.x)
		res, _ := fun.Call(g)
		return res
	}).(Applicative)
}

func (x failure) Chain(f func(Any) Monad) Monad {
	return x
}

func (x success) Chain(f func(Any) Monad) Monad {
	return f(x.x)
}

func (x failure) Map(f Morphism) Functor {
	return x.Bimap(Identity, f).(Functor)
}

func (x success) Map(f Morphism) Functor {
	return x.Bimap(Identity, f).(Functor)
}

func (x failure) Concat(y Semigroup) Semigroup {
	a := y.(Validation)
	b := a.Bimap(
		concatAnyvals(x.x),
		Identity,
	)
	return b.(Semigroup)
}

func (x success) Concat(y Semigroup) Semigroup {
	a := y.(Functor)
	b := a.Map(concatAnyvals(x.x))
	return b.(Semigroup)
}

// Derived

func (x failure) Fold(f Morphism, g Morphism) Any {
	return f(x.x)
}

func (x success) Fold(f Morphism, g Morphism) Any {
	return g(x.x)
}

func (x failure) Bimap(f Morphism, g Morphism) Monad {
	return NewFailure(f(x.x))
}

func (x success) Bimap(f Morphism, g Morphism) Monad {
	return NewSuccess(g(x.x))
}
