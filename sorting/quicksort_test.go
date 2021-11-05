package sorting_test

import (
	"fmt"
	"github.com/denismitr/gds/sorting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuick_Integers(t *testing.T) {
	tt := []struct{
		in []int
	 	exp []int
	}{
		{in: []int{2,3,4,1,9}, exp: []int{1,2,3,4,9}},
		{in: []int{4,1}, exp: []int{1,4}},
		{in: []int{4,9,1}, exp: []int{1,4,9}},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			qs := sorting.NewQuick()
			qs.Integers(tc.in)

			assert.Equal(t, tc.exp, tc.in)
		})
	}
}

func TestQuick_Floats(t *testing.T) {
	tt := []struct{
		in []float64
		exp []float64
	}{
		{in: []float64{2.3,3.9,4.8,1.2,9.99867}, exp: []float64{1.2,2.3,3.9,4.8,9.99867}},
		{in: []float64{4.0,1.3}, exp: []float64{1.3,4.0}},
		{in: []float64{4.4,9.1,1.7}, exp: []float64{1.7,4.4,9.1}},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			qs := sorting.NewQuick()
			qs.Floats(tc.in)

			assert.Equal(t, tc.exp, tc.in)
		})
	}
}

func BenchmarkQuick_Floats(b *testing.B) {
	in := []float64{2.3,3.9,4.8,1.2,9.99867,8,9873.8773, 0.86763, -9, -7.2, 88}

	s := sorting.FloatSlice(in)
	for i := 0; i < b.N; i++ {
		sorting.QuickSort(s)
	}
}
