package junk

type Result(type T) struct {
	Ok T
	Err error
}

type ResultF(S, T) Result(S)

func (rs *ResultFunctor(S, T)) FMap(f func(S) T) Result(T) {
	if rs.Err != nil {
		return Result(T){Err: rs.Err}
	}
	return Result(T){Ok: f(rs.Ok)}
}

contract Functor(S, FS, T, FT) {
	FS FMap(func (S) T) FT
}

func stuff() {
	rone := Result{Ok: 1}
	rfour := Result(int)(ResultF(int, int)(rone).Fmap(func(x int) {return x + 3}))

}


contract Unit(T, MT) {
	MT UnitFrom(T)
}

contract UnitJoin(T, MT, MMT) {
	Unit(T, MT)
	Unit(MT, MMT)
	MMT Join() MT
}

contract Monadic(S, MS, T, MT, MMT) {
	Unit(S, MS)
	UnitJoin(T, MT, MMT)
	Functor(S, MS, T, MT)
	Functor(S, MS, MT, MMT)
}

contract Monadic3(S, MS, T, MT, MMT, U, MU, MMU) {
	Monadic(S, MS, T, MT, MMT)
	Monadic(T, MT, U, MU, MMU)
}

func Bind(type S, MS, T, MT, MMT)(f func(S) MT) func(MS) MT {
	return func(ms) {
		return ms.FMap(f).Join()
	}
}

func Kliesli(type S, MS, T, MT, MMT, U, MU, MMU Monadic3)(g func(T) MU, f func(S) MT) func(S) MU {
	return func(s) {
		return Bind(g)(f(s))
	}
}

func ToResult01(type T)(f func() (T, err)) func() Result(T) {
	return func() {
		t, err := f()
		return Result{t, err}
	}
}

func ToResult11(type S, T)(f func(S) (T, err)) func(S) Result(T) {
	return func(s) {
		t, err := f(s)
		return Result{t, err}
	}
}

func FromResult11(type S, T)(f func(S) Result(T)) func(S) (T, err) {
	return func(s) {
		rt := f(s)
		return rt.Ok, rt.Err
	}
}

func MapResult01(type T)(f func() (T, err)) func() Result(T) {
	return ToResult01(f)
}

func MapResult11(type S, T)(f func(S) (T, err)) func(Result(S)) Result(T) {
	return func(rs) {
		var zero T
		if rs.Err != nil {
			return Result{zero, rs.Err}
		}
		return ToResult11(f)(rs.Ok)
	}
}

func MapResult10(type S, T)(f func(S) err) func(Result(S)) err {
	return func(rs) {
		var zero T
		if rs.Err != nil {
			return rs.Err
		}
		return f(rs.Ok)
	}
}


func hgf() err {
	y, err := f()
	if err != nil {
		return err
	}
	z, err := g(y)
	if err != nil {
		return err
	}
	return h(z)
}

func hgfM() err {
	return MapResult10(h)(MapResult11(g)(MapResult01(f)()))
}
