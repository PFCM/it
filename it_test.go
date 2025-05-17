package it

import (
	"iter"
	"slices"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestZip(t *testing.T) {
	type pair struct {
		A int
		B string
	}
	pairs := func(i iter.Seq2[int, string]) (ps []pair) {
		for a, b := range i {
			ps = append(ps, pair{a, b})
		}
		return
	}

	for _, c := range []struct {
		name string
		as   []int
		bs   []string
		want []pair
	}{{
		name: "equal-sizes",
		as:   []int{1, 2, 3, 4},
		bs:   []string{"a", "b", "c", "d"},
		want: []pair{
			{1, "a"},
			{2, "b"},
			{3, "c"},
			{4, "d"},
		},
	}, {
		name: "first-shorter",
		as:   []int{1, 2, 3},
		bs:   []string{"a", "b", "c", "d"},
		want: []pair{{1, "a"}, {2, "b"}, {3, "c"}},
	}, {
		name: "second-shorter",
		as:   []int{1, 2, 3, 4},
		bs:   []string{"a", "b", "c"},
		want: []pair{{1, "a"}, {2, "b"}, {3, "c"}},
	}} {
		t.Run(c.name, func(t *testing.T) {
			got := pairs(Zip(slices.Values(c.as), slices.Values(c.bs)))
			if d := cmp.Diff(got, c.want); d != "" {
				t.Fatalf("unexpected zip (-got, +want):\n%v", d)
			}
		})
	}
}

func TestChain(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	for i := range 10 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			in := make([]iter.Seq[int], i)
			for j := range in {
				in[j] = slices.Values(values)
			}

			got := slices.Collect(Chain(in...))
			want := slices.Repeat(values, i)
			// Special case, should only happen when i ==
			// 0, turns into a nil vs empty slice
			// situation.
			if got == nil {
				got = []int{}
			}
			if want == nil {
				want = []int{}
			}

			if d := cmp.Diff(got, want); d != "" {
				t.Fatalf("chain mismatch (-got, +want):\n%v", d)
			}
		})
	}
}
