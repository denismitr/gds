package daryheap_test

import (
	"encoding/binary"
	"github.com/denismitr/gds/daryheap"
	"github.com/stretchr/testify/assert"
	"hash/fnv"
	"testing"
)

type simpleTestElement struct {
	v string
}

func (el *simpleTestElement) Hash() uint64 {
	h := fnv.New64a()

	if err := binary.Write(h, binary.LittleEndian, []byte(el.v)); err != nil {
		panic(err)
	}

	return h.Sum64()
}

type testableString string

func (s testableString) Hash() uint64 {
	h := fnv.New64a()

	if err := binary.Write(h, binary.LittleEndian, []byte(s)); err != nil {
		panic(err)
	}

	return h.Sum64()
}

func Test_DHeap_New(t *testing.T) {
	t.Run("it can be initialized empty", func(t *testing.T) {
		dh, err := daryheap.New(2)
		if err != nil {
			t.Fatal(err)
		}

		if !dh.Empty() {
			t.Fatalf("expected DHeap to be empty")
		}
	})

	t.Run("it can be initialized with one element", func(t *testing.T) {
		dh, err := daryheap.New(2)
		if err != nil {
			t.Fatal(err)
		}

		dh.Insert(&simpleTestElement{"foo"}, 1)
		if dh.Empty() {
			t.Fatalf("expected DHeap not to be empty")
		}

		pv, err := dh.Peek()
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		got, want := pv.(*simpleTestElement).v, "foo"
		if got != want {
			t.Fatalf("Invalid dheap.Peek result: Want '%s', got '%s'", want, got)
		}

		tv, err := dh.Top()
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		tGot, tWant := tv.(*simpleTestElement).v, "foo"
		if tGot != tWant {
			t.Fatalf("Invalid dheap.Top result: Want '%s', got '%s'", tWant, tGot)
		}

		if !dh.Empty() {
			t.Fatalf("expected DHeap to be empty after Top()")
		}
	})
}

func TestDHeap_Insert(t *testing.T) {
	t.Run("2-ary max heap", func(t *testing.T) {
		dh, err := daryheap.New(2)
		if err != nil {
			t.Fatal(err)
		}

		dh.Insert(&simpleTestElement{"foo1"}, 20)
		dh.Insert(&simpleTestElement{"foo2"}, 2)
		dh.Insert(&simpleTestElement{"foo3"}, 987)
		dh.Insert(&simpleTestElement{"foo4"}, 454)
		dh.Insert(&simpleTestElement{"foo5"}, -2)
		dh.Insert(&simpleTestElement{"foo6"}, 490)

		assert.True(t, dh.Contains(&simpleTestElement{"foo1"}))

		foo3, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo3", foo3.(*simpleTestElement).v)

		foo6, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo6", foo6.(*simpleTestElement).v)

		foo4, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo4", foo4.(*simpleTestElement).v)

		foo1, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo1", foo1.(*simpleTestElement).v)

		foo2, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo2", foo2.(*simpleTestElement).v)

		foo5, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo5", foo5.(*simpleTestElement).v)

		assert.True(t, dh.Empty())
	})
}

func TestDaryHeap_UpdatePriority(t *testing.T) {
	t.Run("make first element - last", func(t *testing.T) {
		dh, err := daryheap.New(2)
		if err != nil {
			t.Fatal(err)
		}

		dh.Insert(&simpleTestElement{"foo1"}, 20)
		dh.Insert(&simpleTestElement{"foo2"}, 2)
		dh.Insert(&simpleTestElement{"foo3"}, 987)
		dh.Insert(&simpleTestElement{"foo4"}, 454)
		dh.Insert(&simpleTestElement{"foo5"}, -2)
		dh.Insert(&simpleTestElement{"foo6"}, 490)

		foo3, err := dh.Peek()
		if err != nil {
			t.Fatalf("Peek() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo3", foo3.(*simpleTestElement).v)

		if err := dh.UpdatePriority(&simpleTestElement{"foo3"}, -10); err != nil {
			t.Fatalf("could not update priority for foo3: %v", err)
		}

		foo6, err := dh.Top()
		if err != nil {
			t.Fatalf("Top() returned unexpected error %v", err)
		}
		assert.Equal(t, "foo6", foo6.(*simpleTestElement).v)
	})
}

func TestDaryHeap_Remove(t *testing.T) {
	t.Run("remove some of the elements", func(t *testing.T) {
		dh, err := daryheap.New(3)
		if err != nil {
			t.Fatal(err)
		}

		dh.Insert(testableString("foo"), 45.4)
		dh.Insert(testableString("bar"), 46.4)
		dh.Insert(testableString("baz"), 5.3)
		dh.Insert(testableString("abc123"), -45.4)
		dh.Insert(testableString("cba123abc"), 145.9)
		dh.Insert(testableString("foobar"), -148.9)
		dh.Insert(testableString("barBaz"), 95.932)
		dh.Insert(testableString("fooBaz"), -91.932)

		assert.Equal(t, 8, dh.Size())

		// test that some items exist
		assert.True(t, dh.Contains(testableString("barBaz")))
		assert.True(t, dh.Contains(testableString("cba123abc")))
		assert.True(t, dh.Contains(testableString("bar")))
		assert.True(t, dh.Contains(testableString("foobar")))

		if err := dh.Remove(testableString("foobar")); err != nil {
			t.Fatal(err)
		}

		assert.False(t, dh.Contains(testableString("foobar")))

		cba123abc, err := dh.Peek()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testableString("cba123abc"), cba123abc)

		// do: remove the first element
		if err := dh.Remove(cba123abc.(testableString)); err != nil {
			t.Fatal(err)
		}

		// top element should have changed
		barBaz, err := dh.Peek()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testableString("barBaz"), barBaz)
	})
}
