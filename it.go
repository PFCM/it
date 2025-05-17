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
