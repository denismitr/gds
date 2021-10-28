package linkedlist

type Node struct {
	data interface{}
	next *Node
}

type LinkedList struct {
	head *Node
	size int
}

func New() *LinkedList {
	return &LinkedList{
		size: 0,
		head: nil,
	}
}

func (ll *LinkedList) Add(data interface{}) {
	newNode := &Node{data: data, next: nil}

	if ll.head == nil {
		ll.head = newNode
		ll.size = 1
		return
	}

	curr := ll.head
	for curr.next != nil {
		curr = curr.next
	}

	curr.next = newNode
	ll.size++
}

func (ll *LinkedList) Items() []interface{} {
	result := make([]interface{}, 0, ll.size)
	curr := ll.head
	for curr != nil {
		result = append(result, curr.data)
		curr = curr.next
	}
	return result
}

func (ll *LinkedList) Reverse() {
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

func (ll *LinkedList) Count() int {
	return ll.size
}