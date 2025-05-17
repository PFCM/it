package it

import "iter"

// Fold performs a left fold across the iterator using the provided combining
// function and initial value.
func Fold[A, B any](it iter.Seq[A], z B, f func(A, B) B) B {
	b := z
	for a := range it {
		b = f(a, b)
	}
	return b
}

// All is a specialised fold for iterators of bools that returns true iff all of
// the values yielded by the iterator are true.
func All(bs iter.Seq[bool]) bool {
	// Could be:
	// return Fold(bs, true, func(a, b bool) bool {
	// 	return a && b
	// })
	// but may as well stop early if we can.
	for b := range bs {
		if !b {
			return false
		}
	}
	return true
}
