package it

import "iter"

// CollectErr collects all of the A elements from the iterator, up until the
// first non-nil error. When a non-nil error is encountered it is immediately
// returned, along with with any values successfully retrieved up to that point.
// Note that any value returned alongside the error will _not_ be returned; the
// assumption is that the iterator either returns an error or a value and not
// both.
func CollectErr[A any](i iter.Seq2[A, error]) ([]A, error) {
	var values []A
	for a, err := range i {
		if err != nil {
			return values, err
		}
		values = append(values, a)
	}
	return values, nil
}

