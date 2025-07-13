// package it provides some tools for working with iterators, inspired by
// Python's itertools.
package it

import (
	"iter"
	"slices"
)

// Zip returns an iterator that iterates through a and b at the same time,
// yielding pairs of adjacent items. The returned iterator stops as soon as
// either as or bs runs out of items.
func Zip[A, B any](as iter.Seq[A], bs iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextA, stopA := iter.Pull(as)
		defer stopA()
		for b := range bs {
			a, ok := nextA()
			if !ok {
				return
			}
			if !yield(a, b) {
				return
			}
		}
	}
}

// Enumerate returns an iterator that pairs each element in the provided
// sequence with its index in the sequence, starting from 0.
func Enumerate[A any](it iter.Seq[A]) iter.Seq2[int, A] {
	return func(yield func(int, A) bool) {
		j := 0
		for i := range it {
			if !yield(j, i) {
				return
			}
			j++
		}
	}
}

// Chain takes a number of iterators and returns a single iterator that yields
// of the values from all of the iterators in sequence, starting with the first
// argument, then the second and so on.
func Chain[A any](its ...iter.Seq[A]) iter.Seq[A] {
	return Concat(slices.Values(its))
}

// Concat is like chain, but accepts an iterator of iterators.
func Concat[A any](its iter.Seq[iter.Seq[A]]) iter.Seq[A] {
	return func(yield func(A) bool) {
		for it := range its {
			for a := range it {
				if !yield(a) {
					return
				}
			}
		}
	}
}

// Batch returns an iterator that yields batches of n consecutive values from
// the provided iterator. The last batch may be smaller. The yielded slice is
// only valid until the next value is yields (it is reused between batches).
func Batch[A any](i iter.Seq[A], n int) iter.Seq[[]A] {
	return func(yield func([]A) bool) {
		if n == 0 {
			return
		}
		batch := make([]A, 0, n)
		for a := range i {
			batch = append(batch, a)
			if len(batch) == n {
				if !yield(batch) {
					return
				}
				batch = batch[:0]
			}
		}
		if len(batch) > 0 {
			yield(batch)
		}
	}
}

// Limit returns a new iterator that yields the first n values from the provided
// iterator and then stops. If the parent iterator has fewer than n values, the
// returned child iterator will just stop when it runs out.
func Limit[A any](i iter.Seq[A], n int) iter.Seq[A] {
	return func(yield func(A) bool) {
		if n == 0 {
			return
		}
		for i, a := range Enumerate(i) {
			if !yield(a) {
				return
			}
			if i == n-1 {
				return
			}
		}
	}
}

// Map applies a function to every item in the iterator.
func Map[A, B any](as iter.Seq[A], f func(A) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range as {
			if !yield(f(a)) {
				return
			}
		}
	}
}

// Map1x2 maps an iter.Seq to an iter.Seq2 by applying the provided function to
// each item in turn.
func Map1x2[A, B, C any](as iter.Seq[A], f func(A) (B, C)) iter.Seq2[B, C] {
	return func(yield func(B, C) bool) {
		for a := range as {
			if !yield(f(a)) {
				return
			}
		}
	}
}

// Map2x1 maps an iter.Seq2 to an iter.Seq by applying the provided function to
// each pair of items in turn.
func Map2x1[A, B, C any](abs iter.Seq2[A, B], f func(A, B) C) iter.Seq[C] {
	return func(yield func(C) bool) {
		for a, b := range abs {
			if !yield(f(a, b)) {
				return
			}
		}
	}
}

// Map2x2 maps an iter.Seq2 to an iter.Seq2 by applying the provided function to
// each pair of items in turn.
func Map2x2[A, B, C, D any](abs iter.Seq2[A, B], f func(A, B) (C, D)) iter.Seq2[C, D] {
	return func(yield func(C, D) bool) {
		for a, b := range abs {
			if !yield(f(a, b)) {
				return
			}
		}
	}
}

// Const returns an iterator that continually yields the provided value,
// forever. Note that this is an infinite iterator, intended to be used with
// something like Zip or Take that will stop early.
func Const[A any](a A) iter.Seq[A] {
	return func(yield func(A) bool) {
		for yield(a) {
		}
	}
}

// Take returns an iterator that yields at most the first n elements of the
// provided iterator and then stops.
func Take[A any](it iter.Seq[A], n int) iter.Seq[A] {
	return func(yield func(A) bool) {
		i := 0
		for a := range it {
			if i >= n {
				return
			}
			if !yield(a) {
				return
			}
			i++
		}
	}
}

// TakeWhile returns an iterator that yields the (possibly empty) prefix of the
// provided iterator for which the given predicate returns true. The returned
// iterator finishes as soon as it yields a value for which p returns false.
func TakeWhile[A any](it iter.Seq[A], p func(A) bool) iter.Seq[A] {
	return func(yield func(A) bool) {
		for a := range it {
			if !p(a) {
				return
			}
			if !yield(a) {
				return
			}
		}
	}
}

// Filter returns an iterator which yields only those values in it for which p
// returns true.
func Filter[A any](it iter.Seq[A], p func(A) bool) iter.Seq[A] {
	return func(yield func(A) bool) {
		for a := range it {
			if !p(a) {
				continue
			}
			if !yield(a) {
				return
			}
		}
	}
}

// Pair is just a pair of two elements, for occasions where we need to do things
// like collect the values in an iter.Seq2.
type Pair[A, B any] struct {
	A A
	B B
}

// NewPair creates a new pair. It is useful to have this defined as a function,
// such as when implementing Collect2.
func NewPair[A, B any](a A, b B) Pair[A, B] { return Pair[A, B]{A: a, B: b} }

// Values returns the values of the pair.
func (p Pair[A, B]) Values() (A, B) { return p.A, p.B }

// Collect2 is like slices.Collect, but works with iter.Seq2, returning all of
// the results as Pairs.
func Collect2[A, B any](i iter.Seq2[A, B]) []Pair[A, B] {
	return slices.Collect(Map2x1(i, NewPair))
}

// Unpair is a convenience for turning an iter.Seq[Pair[A, B]] into an
// iter.Seq2[A, B].
func Unpair[A, B any](i iter.Seq[Pair[A, B]]) iter.Seq2[A, B] {
	return Map1x2(i, Pair[A, B].Values)
}
