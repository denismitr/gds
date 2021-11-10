package sorting

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

type Sorter interface {
	Integers([]int)
	Floats([]float64)
	Strings([]string)
}

type (
	IntegerSlice []int
	StringSlice  []string
	FloatSlice []float64
)

func (f FloatSlice) Len() int {
	return len(f)
}

func (f FloatSlice) Less(i, j int) bool {
	return f[i] < f[j]
}

func (f FloatSlice) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (s StringSlice) Len() int {
	return len(s)
}

func (s StringSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s StringSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IntegerSlice) Len() int {
	return len(s)
}

func (s IntegerSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s IntegerSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

