package sorting

import "errors"

type Merge struct {}

func NewMerge() Merge {
	return Merge{}
}

func (Merge) Integers(s []int) {
	MergeSort(IntegerSlice(s))
}

func (Merge) Floats(s []float64) {
	MergeSort(FloatSlice(s))
}

func (Merge) Strings(s []string) {
	MergeSort(StringSlice(s))
}

func MergeSort(s Sortable) {
	switch typedSlice := s.(type) {
	case IntegerSlice:
		tmp := make(IntegerSlice, s.Len())
		mergeSortIntegers(typedSlice, tmp, 0, s.Len() - 1)
	case FloatSlice:
		tmp := make(FloatSlice, s.Len())
		mergeSortFloats(typedSlice, tmp, 0, s.Len() - 1)
	default:
		panic(errors.New("invalid slice type"))
	}
}

func mergeSortIntegers(s IntegerSlice, tmp IntegerSlice, left, right int) {
	if left >= right {
		return
	}

	middle := (left + right) / 2

	mergeSortIntegers(s, tmp,  left, middle)
	mergeSortIntegers(s, tmp, middle + 1, right)
	mergeIntegerHalves(s, tmp, left, right)
}

func mergeIntegerHalves(s IntegerSlice, tmp IntegerSlice, left int, right int) {
	middle := (left + right) / 2
	size := right - left + 1

	l := left
	r := middle + 1
	t := left
	for l <= middle && r <= right {
		if s[l] <= s[r] {
			tmp[t] = s[l]
			l++
		}  else {
			tmp[t] = s[r]
			r++
		}
		t++
	}

	leftRemainder := middle - l + 1
	rightRemainder := right - r + 1
	if leftRemainder > 0 {
		copy(tmp[t:t+leftRemainder], s[l:l+leftRemainder])
	} else if rightRemainder > 0 {
		copy(tmp[t:t+rightRemainder], s[r:l+rightRemainder])
	}

	copy(s[left:left+size], tmp[left:left+size])
}

func mergeSortFloats(s FloatSlice, tmp FloatSlice, left, right int) {
	if left >= right {
		return
	}

	middle := (left + right) / 2

	mergeSortFloats(s, tmp,  left, middle)
	mergeSortFloats(s, tmp, middle + 1, right)
	mergeFloatHalves(s, tmp, left, right)
}

func mergeFloatHalves(s FloatSlice, tmp FloatSlice, left int, right int) {
	middle := (left + right) / 2
	size := right - left + 1

	l := left
	r := middle + 1
	t := left
	for l <= middle && r <= right {
		if s[l] <= s[r] {
			tmp[t] = s[l]
			l++
		}  else {
			tmp[t] = s[r]
			r++
		}
		t++
	}

	leftRemainder := middle - l + 1
	rightRemainder := right - r + 1
	if leftRemainder > 0 {
		copy(tmp[t:t+leftRemainder], s[l:l+leftRemainder])
	} else if rightRemainder > 0 {
		copy(tmp[t:t+rightRemainder], s[r:l+rightRemainder])
	}

	copy(s[left:left+size], tmp[left:left+size])
}
