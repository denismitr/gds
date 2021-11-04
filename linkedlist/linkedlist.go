package linkedlist

import "github.com/denismitr/gds/internal/utils"

type Node struct {
	data interface{}
	next *Node
}

type LinkedList struct {
	locker utils.Locker
	head *Node
	size int
}

func New(locks bool) *LinkedList {
	ll := &LinkedList{
		size: 0,
		head: nil,
	}

	if locks {
		ll.locker = &utils.MutexLock{}
	} else {
		ll.locker = utils.NullLocker{}
	}

	return ll
}

// Append any data to linked list and returns an index
func (ll *LinkedList) Append(data interface{}) int {
	ll.locker.WriteLock()
	defer ll.locker.WriteUnlock()

	newNode := &Node{data: data, next: nil}

	if ll.head == nil {
		ll.head = newNode
		ll.size = 1
		return 0
	}

	curr := ll.head
	for curr.next != nil {
		curr = curr.next
	}

	curr.next = newNode
	ll.size++
	return ll.size - 1
}

// Slice returns data stored at linked list nodes as a slice
func (ll *LinkedList) Slice() []interface{} {
	ll.locker.ReadLock()
	defer ll.locker.ReadUnlock()

	result := make([]interface{}, 0, ll.size)
	curr := ll.head
	for curr != nil {
		result = append(result, curr.data)
		curr = curr.next
	}
	return result
}

// Reverse the linked list
func (ll *LinkedList) Reverse() {
	ll.locker.WriteLock()
	defer ll.locker.WriteUnlock()

	if ll.head == nil {
		return
	}

	curr := ll.head
	var prev *Node
	var next *Node

	for curr != nil {
		next = curr.next
		curr.next = prev
		prev = curr
		curr = next
	}

	ll.head = prev
}

// Pop takes data at head, returns it and sets head to head.next
func (ll *LinkedList) Pop() (interface{}, bool) {
	ll.locker.WriteLock()
	defer ll.locker.WriteUnlock()

	if ll.head == nil {
		return nil, false
	}

	n := ll.head
	ll.head = ll.head.next
	ll.size -= 1

	return n.data, true
}

func (ll *LinkedList) PeakAt(index int) (interface{}, bool) {
	ll.locker.ReadLock()
	defer ll.locker.ReadUnlock()

	if ll.head == nil || index > ll.size - 1 {
		return nil, false
	}

	curr := ll.head
	for i := 1; i <= index; i++ {
		curr = curr.next
	}

	return curr.data, true
}

// Size of the linked list
func (ll *LinkedList) Size() int {
	ll.locker.ReadLock()
	defer ll.locker.ReadUnlock()
	return ll.size
}

func (ll *LinkedList) Empty() bool {
	ll.locker.ReadLock()
	defer ll.locker.ReadUnlock()
	return ll.size == 0
}