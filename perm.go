package it

import "iter"

// Perm returns an iterator that yields all permutations of the provided slice.
// It shuffles the objects in place, and always yields the same slice, so care
// must be taken if re-using the yielded values.
func Perm[E any, S ~[]E](data S) iter.Seq[S] {
	// According to the benchmarks, the iterative version is a bit quicker,
	// at the cost of one extra alloc. Which is a reasonable price to pay
	// given callers are unlikely to be trying to permute anything too huge.
	return permIter(data)
}

func permRec[E any, S ~[]E](data S) iter.Seq[S] {
	return func(yield func(S) bool) {
		switch len(data) {
		case 0:
			return
		case 1:
			yield(data)
			return
		}
		ret := make(S, len(data))
		// Heap's algorithm.
		var perm func(k int, s S) bool
		perm = func(k int, s S) bool {
			if k == 1 {
				copy(ret, s)
				return yield(ret)
			}
			if !perm(k-1, s) {
				return false
			}
			for i := range k - 1 {
				if k%2 == 0 {
					s[i], s[k-1] = s[k-1], s[i]
				} else {
					s[0], s[k-1] = s[k-1], s[0]
				}
				if !perm(k-1, s) {
					return false
				}
			}
			return true
		}
		perm(len(data), data)
	}
}

func permIter[E any, S ~[]E](data S) iter.Seq[S] {
	return func(yield func(S) bool) {
		if len(data) == 0 {
			return
		}
		ret := make(S, len(data))
		yieldCopy := func() bool {
			copy(ret, data)
			return yield(ret)
		}
		// via https://sedgewick.io/wp-content/uploads/2022/03/2002PermGeneration.pdf
		c := make([]int, len(data))
		for i := range c {
			c[i] = 0
		}
		if !yieldCopy() {
			return
		}

		for i := 0; i < len(data); {
			if c[i] < i {
				if i%2 == 0 {
					data[0], data[i] = data[i], data[0]
				} else {
					data[c[i]], data[i] = data[i], data[c[i]]
				}
				if !yieldCopy() {
					return
				}
				c[i] += 1
				i = 1
			} else {
				c[i] = 0
				i += 1
			}
		}
	}
}
