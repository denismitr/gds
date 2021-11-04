package linkedlist_test

import (
	"fmt"
	"github.com/denismitr/gds/linkedlist"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestLinkedList_NoConcurrency(t *testing.T) {
	tt := []struct{
		locks bool
		seed []interface{}
		reversed []interface{}
	}{
		{
			seed: []interface{}{1,2,3,4,5},
			reversed: []interface{}{5,4,3,2,1},
		},
		{
			locks: true,
			seed: []interface{}{1,2,3,4,5,43554,34535,1234,98775,7,0,99,34,80},
			reversed: []interface{}{80,34,99,0,7,98775,1234,34535,43554,5,4,3,2,1},
		},
		{
			seed: []interface{}{"a","b","c","d","e"},
			reversed: []interface{}{"e","d","c","b","a"},
		},
		{
			seed: []interface{}{1},
			reversed: []interface{}{1},
		},
		{
			seed: []interface{}{},
			reversed: []interface{}{},
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			ll := linkedlist.New(tc.locks)
			assert.True(t, ll.Empty())

			for _, data := range tc.seed {
				idx := ll.Append(data)
				d, ok := ll.PeakAt(idx)
				assert.Truef(t, ok, "expected to get data at index %d", idx)
				assert.Equalf(t, tc.seed[idx], d, "expected to get data at index %d", idx)
			}

			actualSize, expectedSize := ll.Size(), len(tc.seed)
			if actualSize != expectedSize {
				t.Errorf("expected Linked List size to be %d, got %d", expectedSize, actualSize)
			}

			if expectedSize > 0 {
				assert.False(t, ll.Empty())
			}

			// let's check the contents
			for idx := range tc.seed {
				d, ok := ll.PeakAt(idx)
				assert.Truef(t, ok, "expected to get data at index %d", idx)
				assert.Equalf(t, tc.seed[idx], d, "expected to get data at index %d", idx)
			}

			items := ll.Slice()
			if len(items) != expectedSize {
				t.Errorf("expected Linked List items count to be %d, got %d", expectedSize, len(items))
			}

			assertSlicesEqual(t, tc.seed, items)

			ll.Reverse()
			reversedItems := ll.Slice()
			if len(reversedItems) != expectedSize {
				t.Errorf("expected Linked List items count to be %d, got %d", expectedSize, len(items))
			}

			assertSlicesEqual(t, tc.reversed, reversedItems)

			for j := len(tc.seed) - 1; j >= 0; j-- {
				d, ok := ll.Pop()
				assert.Truef(t, ok, "expected data to be popped with success")
				assert.Equalf(t, d, tc.seed[j], "expected popped data to be equal to seeded item")
			}

			d, ok := ll.Pop()
			assert.False(t, ok, "there should be nothing left to pop")
			assert.Nilf(t, d, "there should be nothing left to pop")
		})
	}
}

func TestLinkedList_WithConcurrency(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const min = 10
	const max = 300
	const concurrency = 10
	const size = 150

	t.Run("append, peak and pop concurrently", func(t *testing.T) {
		ll := linkedlist.New(true)

		var wg sync.WaitGroup

		for i := 0; i < concurrency; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for j := 0; j < size; j++ {
					d := rand.Intn(max - min + 1) + min
					idx := ll.Append(d)
					assert.Greater(t, ll.Size(), 0)
					d2, ok := ll.PeakAt(idx)
					assert.True(t, ok)
					assert.Equalf(t, d, d2, "expected appended and peaked values to be equal")
				}
			}()
		}

		wg.Wait()

		for i := 0; i < concurrency; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for j := 0; j < size; j++ {
					assert.Greater(t, ll.Size(), 0)
					assert.Falsef(t, ll.Empty(), "linked list should not be empty")

					d, ok := ll.Pop()
					assert.Truef(t, ok, "linked list should not be empty")

					n := d.(int)
					assert.GreaterOrEqual(t, n, min)
					assert.LessOrEqual(t, n, max)
				}
			}()
		}

		wg.Wait()
	})
}

func assertSlicesEqual(t *testing.T, a, b []interface{}) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatalf("expected two slices too have equal size, got %d != %d", len(a), len(b))
	}

	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			t.Fatalf("expected items of twol slices at index %d, got %+v != %+v", i, a[i], b[i])
		}
	}
}