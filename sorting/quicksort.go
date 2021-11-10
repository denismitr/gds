package sorting

import (
	"math/rand"
)

func QuickSort(s Sortable) {
	quicksort(s, 0, s.Len() - 1)
}

type Quick struct {}

func (Quick) Integers(s []int) {
	QuickSort(IntegerSlice(s))
}

func (Quick) Floats(s []float64) {
	QuickSort(FloatSlice(s))
}

func (Quick) Strings(s []string) {
	QuickSort(StringSlice(s))
}

func NewQuick() Sorter {
	return Quick{}
}

func quicksort(s Sortable, left, right int) {
	if right <= left {
		return
	}

	pivot := partition(s, left, right)
	quicksort(s, left, pivot - 1)
	quicksort(s, pivot + 1, right)
}

func partition(s Sortable, left, right int) int {
	pivot := rand.Intn(right - left) + left
	s.Swap(pivot, right)

	for j := left; j <= right - 1; j++ {
		// If current element is smaller than the pivot
		if s.Less(j, right) {
			s.Swap(left, j)
			left++
		}
	}

	s.Swap(left, right)
	return left
}
