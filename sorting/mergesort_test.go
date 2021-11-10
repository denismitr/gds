package sorting_test

import (
	"fmt"
	"github.com/denismitr/gds/sorting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeSort_Integers(t *testing.T) {
	tt := []struct{
		in []int
		exp []int
	}{
		{in: []int{2,3,4,1,9}, exp: []int{1,2,3,4,9}},
		{in: []int{2,3,4,1,9,89,-876,414,0}, exp: []int{-876,0,1,2,3,4,9,89,414}},
		{in: []int{4,1}, exp: []int{1,4}},
		{in: []int{4,9,1}, exp: []int{1,4,9}},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			qs := sorting.NewMerge()
			qs.Integers(tc.in)

			assert.Equal(t, tc.exp, tc.in)
		})
	}
}

func TestMergeSort_Floats(t *testing.T) {
	tt := []struct{
		in []float64
		exp []float64
	}{
		{in: []float64{2.1,3.4,4.2,1.1,9.9}, exp: []float64{1.1,2.1,3.4,4.2,9.9}},
		{in: []float64{2.9,3.01,4.1,1,9.89,89.333,-876.8,414,0}, exp: []float64{-876.8,0,1,2.9,3.01,4.1,9.89,89.333,414}},
		{in: []float64{4.2009,1.9999}, exp: []float64{1.9999,4.2009}},
		{in: []float64{4,9,1}, exp: []float64{1,4,9}},
		{in: []float64{}, exp: []float64{}},
		{in: []float64{0.54}, exp: []float64{0.54}},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			qs := sorting.NewMerge()
			qs.Floats(tc.in)

			assert.Equal(t, tc.exp, tc.in)
		})
	}
}
