package it

import (
	"fmt"
	"iter"
	"slices"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPermExplicit(t *testing.T) {
	in := []int{1, 2, 3}
	want := [][]int{
		{1, 2, 3},
		{2, 1, 3},
		{3, 1, 2},
		{1, 3, 2},
		{2, 3, 1},
		{3, 2, 1},
	}

	var got [][]int
	for g := range Perm(in) {
		got = append(got, slices.Clone(g))
	}

	if d := cmp.Diff(got, want); d != "" {
		t.Fatalf("mismatch (-got, +want):\n%v", d)
	}
}

func TestPerms(t *testing.T) {
	for size := range 10 {
		for _, c := range []struct {
			name string
			f    func([]int) iter.Seq[[]int]
		}{
			{name: "recursive", f: permRec[int, []int]},
			{name: "iterative", f: permIter[int, []int]},
		} {
			t.Run(fmt.Sprintf("%d/%s", size, c.name), func(t *testing.T) {
				data := make([]int, size)
				for i := range data {
					data[i] = i
				}
				// If we prove that the function:
				// - returns the right number of permutations, size!
				// - each permutation has the right elements in it
				// - all the permutations are unique,
				//
				// then the it must be correct.
				got := make(map[string]bool)
				for p := range c.f(data) {
					s := fmt.Sprint(p)
					if got[s] {
						t.Fatalf("permutation seen twice:\n%v", s)
					}
					got[s] = true

					if l := len(data); l != size {
						t.Fatalf("invalid permuation: want length %d, got length %d", size, l)
					}
					sort.Ints(p)
					for i, j := range p {
						if i != j {
							t.Fatalf("unexpected element in partition (%d): %v", i, p)
						}
					}
				}
				want := 1
				for i := 1; i < size; i++ {
					want += want * i
				}
				if size == 0 {
					want = 0
				}
				if l := len(got); l != want {
					t.Fatalf("unexpected number of permutations for slice of size %d: want %d, got %d", size, want, l)
				}
			})
		}
	}
}

func BenchmarkPerm(b *testing.B) {
	for size := range 10 {
		data := make([]int, size)
		for i := range data {
			data[i] = i
		}
		for _, c := range []struct {
			name string
			f    func([]int) iter.Seq[[]int]
		}{
			{name: "recursive", f: permRec[int, []int]},
			{name: "iterative", f: permIter[int, []int]},
		} {
			b.Run(fmt.Sprintf("%d/%s", size, c.name), func(b *testing.B) {
				for b.Loop() {
					var x []int
					for p := range c.f(data) {
						x = p
					}
					_ = x
				}
			})
		}
	}
}
