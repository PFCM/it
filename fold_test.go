package it

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFold(t *testing.T) {
	values := []int{1, 2, 3, 4}
	got := Fold(slices.Values(values), []int{}, func(a int, b []int) []int {
		return append(b, a)
	})

	if d := cmp.Diff(got, values); d != "" {
		t.Fatalf("unexpected result (-got, +want): %v", d)
	}
}

func TestAll(t *testing.T) {
	for _, c := range []struct {
		name string
		in   []bool
		want bool
	}{{
		name: "three-true",
		in:   []bool{true, true, true},
		want: true,
	}, {
		name: "first-false",
		in:   []bool{false, true, true},
		want: false,
	}, {
		name: "middle-false",
		in:   []bool{true, false, true},
		want: false,
	}, {
		name: "last-false",
		in:   []bool{true, true, false},
		want: false,
	}, {
		name: "all-false",
		in:   []bool{false, false, false},
		want: false,
	}} {
		t.Run(c.name, func(t *testing.T) {
			if All(slices.Values(c.in)) != c.want {
				t.Fatalf("unexpected result: got %v, want %v\n(input: %v)", !c.want, c.want, c.in)
			}
		})
	}
}
