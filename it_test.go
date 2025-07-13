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

func TestBatch(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	for _, c := range []struct {
		n    int
		want [][]int
	}{{
		n:    0,
		want: nil,
	}, {
		n:    1,
		want: [][]int{{1}, {2}, {3}, {4}, {5}},
	}, {
		n:    2,
		want: [][]int{{1, 2}, {3, 4}, {5}},
	}, {
		n:    3,
		want: [][]int{{1, 2, 3}, {4, 5}},
	}, {
		n:    4,
		want: [][]int{{1, 2, 3, 4}, {5}},
	}, {
		n:    5,
		want: [][]int{{1, 2, 3, 4, 5}},
	}, {
		n:    6,
		want: [][]int{{1, 2, 3, 4, 5}},
	}} {
		t.Run(strconv.Itoa(c.n), func(t *testing.T) {
			// Don't just use slices.Collect, Batch might re-use the
			// slices.
			var got [][]int
			for b := range Batch(slices.Values(in), c.n) {
				got = append(got, slices.Clone(b))
			}

			if d := cmp.Diff(got, c.want); d != "" {
				t.Fatalf("mismatch (-got, +want):\n%v", d)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	for i := range len(data) + 2 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := slices.Collect(Limit(slices.Values(data), i))
			want := data[:min(len(data), i)]

			if len(got) == 0 {
				got = []int{}
			}

			if d := cmp.Diff(got, want); d != "" {
				t.Fatalf("mismatch (-got, +want):\n%v", d)
			}
		})
	}
}

func TestMap(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	out := []string{"1", "2", "3", "4", "5", "6"}

	got := slices.Collect(Map(slices.Values(in), strconv.Itoa))

	if d := cmp.Diff(got, out); d != "" {
		t.Fatalf("mismatch (-got, +want):\n%v", d)
	}
}

func TestMap1x2(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	out := []Pair[int, string]{
		{2, "1"},
		{4, "2"},
		{6, "3"},
		{8, "4"},
		{10, "5"},
		{12, "6"},
	}

	got := Collect2(Map1x2(slices.Values(in), func(x int) (int, string) {
		return x * 2, strconv.Itoa(x)
	}))
	if d := cmp.Diff(got, out); d != "" {
		t.Fatalf("mismatch (-got, +want):\n%v", d)
	}
}

func TestMap2x1(t *testing.T) {
	const n = 10
	in := func(yield func(int, string) bool) {
		for i := range n {
			if !yield(i, strconv.Itoa(i)) {
				return
			}
		}
	}
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	got := slices.Collect(Map2x1(in, func(i int, _ string) int { return i }))
	if d := cmp.Diff(got, want); d != "" {
		t.Fatalf("mismatch (-got, +want):\n%v", d)
	}
}

func TestMap2x2(t *testing.T) {
	in := Unpair(slices.Values([]Pair[int, string]{
		{1, "1"},
		{2, "2"},
		{3, "3"},
		{4, "4"},
		{5, "5"},
	}))
	want := []Pair[string, int]{
		{"1", 1},
		{"2", 2},
		{"3", 3},
		{"4", 4},
		{"5", 5},
	}

	got := Collect2(Map2x2(in, func(i int, s string) (string, int) {
		return s, i
	}))
	if d := cmp.Diff(got, want); d != "" {
		t.Fatalf("mismatch (-got, +want):\n%v", d)
	}
}

func TestTake(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, c := range []struct {
		n    int
		want []int
	}{{
		n:    1,
		want: []int{1},
	}, {
		n:    3,
		want: []int{1, 2, 3},
	}, {
		n:    len(values),
		want: values,
	}, {
		n:    len(values) + 1,
		want: values,
	}, {
		n:    len(values) + 2,
		want: values,
	}, {
		n:    0,
		want: nil,
	}, {
		n:    -1,
		want: nil,
	}} {
		t.Run(strconv.Itoa(c.n), func(t *testing.T) {
			got := slices.Collect(Take(slices.Values(values), c.n))
			if d := cmp.Diff(got, c.want); d != "" {
				t.Fatalf("unexpected result (-got, +want):\n%v", d)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, c := range []struct {
		name string
		p    func(int) bool
		want []int
	}{{
		name: "all",
		p:    func(int) bool { return true },
		want: values,
	}, {
		name: "<4",
		p:    func(i int) bool { return i < 4 },
		want: []int{1, 2, 3},
	}, {
		name: ">4",
		p:    func(i int) bool { return i > 4 },
		want: nil,
	}, {
		name: "odd",
		p:    func(i int) bool { return i%2 == 1 },
		want: []int{1},
	}} {
		t.Run(c.name, func(t *testing.T) {
			got := slices.Collect(TakeWhile(slices.Values(values), c.p))
			if d := cmp.Diff(got, c.want); d != "" {
				t.Fatalf("unexpected result (-got, +want):\n%v", d)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, c := range []struct {
		name string
		p    func(int) bool
		want []int
	}{{
		name: "all",
		p:    func(int) bool { return true },
		want: values,
	}, {
		name: "<4",
		p:    func(i int) bool { return i < 4 },
		want: []int{1, 2, 3},
	}, {
		name: ">4",
		p:    func(i int) bool { return i > 4 },
		want: []int{5, 6, 7, 8, 9, 10},
	}, {
		name: "odd",
		p:    func(i int) bool { return i%2 == 1 },
		want: []int{1, 3, 5, 7, 9},
	}} {
		t.Run(c.name, func(t *testing.T) {
			got := slices.Collect(Filter(slices.Values(values), c.p))
			if d := cmp.Diff(got, c.want); d != "" {
				t.Fatalf("unexpected result (-got, +want):\n%v", d)
			}
		})
	}
}
