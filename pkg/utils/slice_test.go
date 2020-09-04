package utils

import (
	"reflect"
	"testing"
)

func TestStringSliceReflectEqual(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"q", "w", "e", "r", "t"}, []string{"q", "w", "a", "s", "z", "x"}},
	}
	for _, c := range cases {
		result := StringSliceReflectEqual(c.in, c.want)
		if !result {
			t.Errorf("StringSliceReflectEqual(%q) == %q", c.in, c.want)
		}
	}

}

func BenchmarkDeepEqual(b *testing.B) {
	sa := []string{"q", "w", "e", "r", "t"}
	sb := []string{"q", "w", "a", "s", "z", "x"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		StringSliceReflectEqual(sa, sb)
	}
}

func TestStringSliceEqual(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"q", "w", "e", "r", "t"}, []string{"q", "w", "a", "s", "z", "x"}},
	}
	for _, c := range cases {
		result := StringSliceEqual(c.in, c.want)
		if !result {
			t.Errorf("StringSliceEqual(%q) == %q", c.in, c.want)
		}
	}
}

func BenchmarkEqual(b *testing.B) {
	sa := []string{"q", "w", "e", "r", "t"}
	sb := []string{"q", "w", "a", "s", "z", "x"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		StringSliceEqual(sa, sb)
	}
}

func TestUint64SliceReverse(t *testing.T) {
	cases := []struct {
		in, want []uint64
	}{
		{[]uint64{1, 2, 3, 4, 5}, []uint64{5, 4, 3, 2, 1}},
	}
	for _, c := range cases {
		got := Uint64SliceReverse(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("SliceReverseUint64(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
