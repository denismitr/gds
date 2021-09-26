package dary

import (
	"errors"
	"fmt"
	"math"
)

var ErrEmptyHeap = errors.New("empty heap")
var ErrInvalidBranchingFactor = errors.New("branching factor must be greater than 1")
var errCurrentNodeHasNoKid = errors.New("current node has no kid")

type Heap interface {
}

type Identity interface {
	Hash() uint64
}

type element struct {
	value    interface{}
	identity uint64
	priority float64
}

type DHeap struct {
	elements        []element
	identityMap     map[uint64]int
	branchingFactor int
}

func New(branchingFactor int) (*DHeap, error) {
	if branchingFactor < 2 {
		return nil, ErrInvalidBranchingFactor
	}

	h := DHeap{
		elements:        nil,
		identityMap:     make(map[uint64]int),
		branchingFactor: branchingFactor,
	}

	return &h, nil
}

func (h *DHeap) getParentIndex(childIndex int) int {
	if childIndex == 0 {
		return 0
	}

	// for a heap with branching factor D and array’s indices starting at zero
	// parentIndex for node with childIndex would be (childIndex-1) / D
	return (childIndex - 1) / h.branchingFactor
}

func (h *DHeap) bubbleUp(index int) {
	if len(h.elements) < index+1 {
		panic("nodes length and index position mismatch")
	}

	elem := h.elements[index]
	for index > 0 {
		parentIndex := h.getParentIndex(index)
		parent := h.elements[parentIndex]
		if elem.priority > parent.priority {
			h.elements[index] = parent
			h.identityMap[parent.identity] = index
			index = parentIndex
		} else {
			break
		}
	}

	h.elements[index] = elem
	h.identityMap[elem.identity] = index
}

func (h *DHeap) firstLeafIndex() int {
	return (len(h.elements) - 2) / (h.branchingFactor + 1)
}

func (h *DHeap) firstKidIndexOf(index int) int {
	return index*h.branchingFactor + 1
}

// Finds, among the kids of a d-ary heap node, the one child with highest priority.
// In case multiple kids have the same priority, the left-most kid is returned.
// Returns the index of the kid of the current heap node with highest priority,
// or error if current node has no kid.
func (h *DHeap) highestPriorityKidIndex(index int) (int, error) {
	fki := h.firstKidIndexOf(index)
	hSize := len(h.elements)
	if fki > hSize {
		return 0, errCurrentNodeHasNoKid
	}

	lastIndex := minInt(fki, hSize)
	highestPriority := -math.MaxFloat64
	result := fki
	for i := range h.elements[fki:lastIndex] {
		if h.elements[i].priority > highestPriority {
			highestPriority = h.elements[i].priority
			result = i
		}
	}

	return result, nil
}

func (h *DHeap) firstKidIndex(index int) int {
	return h.branchingFactor * index + 1;
}

func (h *DHeap) pushDown(index int) {
	if index < 0 || index >= len(h.elements) {
		panic(fmt.Sprintf("index %d is out of allowed range", index))
	}

	currIndex := index
	elem := h.elements[index]
	smallestKidIndex := h.firstKidIndexOf(index)
	for smallestKidIndex < len(h.elements) {
		lastKidIndexGuard := minInt(h.firstKidIndexOf(index) + h.branchingFactor, len(h.elements))

		for kidIndex := smallestKidIndex; kidIndex < lastKidIndexGuard; kidIndex++ {
			if h.elements[kidIndex].priority > h.elements[smallestKidIndex].priority {
				smallestKidIndex = kidIndex
			}
		}

		//kidIndex, err := h.highestPriorityKidIndex(currIndex)
		//if err != nil {
		//	panic(err) /// ????
		//}

		kid := h.elements[smallestKidIndex]

		if kid.priority > elem.priority {
			h.elements[currIndex] = kid
			h.identityMap[kid.identity] = currIndex
			currIndex = smallestKidIndex
			smallestKidIndex = h.firstKidIndex(currIndex)
		} else {
			break
		}
	}

	h.elements[currIndex] = elem
	h.identityMap[elem.identity] = currIndex
}

func (h *DHeap) heapify() {
	lastInnerElementIndex := h.firstLeafIndex() - 1
	for index := lastInnerElementIndex; index > 0; index-- {
		h.pushDown(index)
	}
}

func (h *DHeap) Insert(v Identity, priority float64) {
	elem := element{value: v, identity: v.Hash(), priority: priority}
	h.elements = append(h.elements, elem)
	h.bubbleUp(len(h.elements) - 1)
}

func (h *DHeap) popValue() interface{} {
	elem := h.elements[0]
	h.elements = append(h.elements[:0], h.elements[1:]...)
	delete(h.identityMap, elem.identity)
	return elem.value
}

func (h *DHeap) Empty() bool {
	return len(h.elements) == 0
}

func (h *DHeap) Top() (interface{}, error) {
	if h.Empty() {
		return nil, ErrEmptyHeap
	}

	if len(h.elements) == 1 {
		return h.popValue(), nil
	}

	elem := h.popValue()
	h.pushDown(0)
	return elem, nil
}

func (h *DHeap) Peek() (interface{}, error) {
	if h.Empty() {
		return nil, ErrEmptyHeap
	}

	return h.elements[0].value, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
