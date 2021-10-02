package daryheap

import (
	"fmt"
	"github.com/denismitr/gds/contracts"
	"github.com/denismitr/gds/internal/utils"
	"github.com/pkg/errors"
	"math"
)

var ErrEmptyHeap = errors.New("empty heap")
var ErrInvalidBranchingFactor = errors.New("branching factor must be greater than 1")
var ErrElementNotFound = errors.New("element not found in identity map")
var errCurrentNodeHasNoKid = errors.New("current node has no kid")

const (
	MinBranchingFactor = 2
	MaxBranchingFactor = 10
)


type options struct {
	isMinHeap bool
	useMutex  bool
}

type OptionFunc func(*options)

type element struct {
	value    interface{}
	identity uint64
	priority float64
}

type priorityComparator func(p1, p2 float64) bool

func minHeapGreaterPriority(p1, p2 float64) bool {
	return p1 < p2
}

func maxHeapGreaterPriority(p1, p2 float64) bool {
	return p1 > p2
}

type DaryHeap struct {
	elements           []element
	identityMap        map[uint64]int
	branchingFactor    int
	locker             utils.Locker
	hasGreaterPriority priorityComparator
}


func New(branchingFactor int, ofs ...OptionFunc) (*DaryHeap, error) {
	opts := options{
		isMinHeap: false,
		useMutex:  true,
	}

	for _, opt := range ofs {
		opt(&opts)
	}

	if branchingFactor < MinBranchingFactor || branchingFactor > MaxBranchingFactor {
		return nil, errors.Wrapf(
			ErrInvalidBranchingFactor,
			"must be between %d and %d",
			MinBranchingFactor,
			MaxBranchingFactor,
		)
	}

	h := DaryHeap{
		elements:        nil,
		identityMap:     make(map[uint64]int),
		branchingFactor: branchingFactor,
	}

	if opts.useMutex {
		h.locker = &utils.MutexLock{}
	} else {
		h.locker = &utils.NullLocker{}
	}

	if opts.isMinHeap {
		h.hasGreaterPriority = minHeapGreaterPriority
	} else {
		h.hasGreaterPriority = maxHeapGreaterPriority
	}

	return &h, nil
}

func (h *DaryHeap) Insert(v contracts.Identity, priority float64) {
	h.locker.WriteLock()
	defer h.locker.WriteUnlock()

	elem := element{value: v, identity: v.Hash(), priority: priority}
	h.elements = append(h.elements, elem)
	h.bubbleUp(len(h.elements) - 1)
}

func (h *DaryHeap) Empty() bool {
	h.locker.ReadLock()
	defer h.locker.ReadUnlock()

	return len(h.elements) == 0
}

func (h *DaryHeap) Size() int {
	h.locker.ReadLock()
	defer h.locker.ReadUnlock()

	return len(h.elements)
}

func (h *DaryHeap) Top() (interface{}, error) {
	h.locker.WriteLock()
	defer h.locker.WriteUnlock()

	if len(h.elements) == 0 {
		return nil, ErrEmptyHeap
	}

	if len(h.elements) == 1 {
		return h.popValue(), nil
	}

	elem := h.popValue()
	h.pushDown(0)
	return elem, nil
}

func (h *DaryHeap) Contains(elem contracts.Identity) bool {
	h.locker.ReadLock()
	defer h.locker.ReadUnlock()
	return h.contains(elem)
}

func (h *DaryHeap) Peek() (interface{}, error) {
	h.locker.ReadLock()
	defer h.locker.ReadUnlock()

	if len(h.elements) == 0 {
		return nil, ErrEmptyHeap
	}

	return h.elements[0].value, nil
}

func (h *DaryHeap) UpdatePriority(elem contracts.Identity, newPriority float64) error {
	h.locker.WriteLock()
	defer h.locker.WriteUnlock()

	index, err := h.findIndexOf(elem)
	if err != nil {
		return err
	}

	oldPriority := h.elements[index].priority
	h.elements[index].priority = newPriority
	if h.hasGreaterPriority(newPriority, oldPriority) {
		h.bubbleUp(index)
	} else if newPriority < oldPriority {
		h.pushDown(index)
	}

	return nil
}

func (h *DaryHeap) Remove(elem contracts.Identity) error {
	h.locker.WriteLock()
	defer h.locker.WriteUnlock()

	index, err := h.findIndexOf(elem)
	if err != nil {
		return err
	}

	n := len(h.elements)
	h.remove(index)
	if n > 1 {
		h.pushDown(index)
	}

	return nil
}

func (h *DaryHeap) contains(elem contracts.Identity) bool {
	_, ok := h.identityMap[elem.Hash()]
	return ok
}

func (h *DaryHeap) getParentIndex(childIndex int) int {
	if childIndex == 0 {
		return 0
	}

	// for a heap with branching factor D and array’s indices starting at zero
	// parentIndex for node with childIndex would be (childIndex-1) / D
	return (childIndex - 1) / h.branchingFactor
}

func (h *DaryHeap) bubbleUp(index int) {
	if len(h.elements) < index+1 {
		panic("nodes length and index position mismatch")
	}

	elem := h.elements[index]
	for index > 0 {
		parentIndex := h.getParentIndex(index)
		parent := h.elements[parentIndex]
		if h.hasGreaterPriority(elem.priority, parent.priority) {
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

func (h *DaryHeap) firstLeafIndex() int {
	return (len(h.elements) - 2) / (h.branchingFactor + 1)
}

func (h *DaryHeap) firstKidIndexOf(index int) int {
	return index*h.branchingFactor + 1
}

// Finds, among the kids of a d-ary heap node, the one child with highest priority.
// In case multiple kids have the same priority, the left-most kid is returned.
// Returns the index of the kid of the current heap node with highest priority,
// or error if current node has no kid.
func (h *DaryHeap) highestPriorityKidIndex(index int) (int, error) {
	fki := h.firstKidIndexOf(index)
	hSize := len(h.elements)
	if fki > hSize {
		return 0, errCurrentNodeHasNoKid
	}

	lastIndex := minInt(fki, hSize)
	highestPriority := -math.MaxFloat64
	result := fki
	for i := range h.elements[fki:lastIndex] {
		if h.hasGreaterPriority(h.elements[i].priority, highestPriority) {
			highestPriority = h.elements[i].priority
			result = i
		}
	}

	return result, nil
}

func (h *DaryHeap) firstKidIndex(index int) int {
	return h.branchingFactor*index + 1
}

func (h *DaryHeap) pushDown(index int) {
	if index < 0 || index >= len(h.elements) {
		panic(fmt.Sprintf("index %d is out of allowed range", index))
	}

	currIndex := index
	elem := h.elements[index]
	smallestKidIndex := h.firstKidIndexOf(index)
	for smallestKidIndex < len(h.elements) {
		lastKidIndexGuard := minInt(h.firstKidIndexOf(index)+h.branchingFactor, len(h.elements))

		for kidIndex := smallestKidIndex; kidIndex < lastKidIndexGuard; kidIndex++ {
			if h.hasGreaterPriority(h.elements[kidIndex].priority, h.elements[smallestKidIndex].priority) {
				smallestKidIndex = kidIndex
			}
		}

		//kidIndex, err := h.highestPriorityKidIndex(currIndex)
		//if err != nil {
		//	panic(err) /// ????
		//}

		kid := h.elements[smallestKidIndex]

		if h.hasGreaterPriority(kid.priority, elem.priority) {
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

func (h *DaryHeap) heapify() {
	lastInnerElementIndex := h.firstLeafIndex() - 1
	for index := lastInnerElementIndex; index > 0; index-- {
		h.pushDown(index)
	}
}

func (h *DaryHeap) popValue() interface{} {
	elem := h.elements[0]
	h.remove(0)
	return elem.value
}

func (h *DaryHeap) remove(index int) {
	delete(h.identityMap, h.elements[index].identity)
	h.elements = append(h.elements[:index], h.elements[index+1:]...)
}

func (h *DaryHeap) findIndexOf(elem contracts.Identity) (int, error) {
	if len(h.elements) == 0 {
		return 0, ErrEmptyHeap
	}

	index, ok := h.identityMap[elem.Hash()]
	if !ok {
		return 0, ErrElementNotFound
	}

	return index, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
