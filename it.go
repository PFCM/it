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
