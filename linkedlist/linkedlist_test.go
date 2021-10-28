package linkedlist_test

import (
	"fmt"
	"github.com/denismitr/gds/linkedlist"
	"reflect"
	"testing"
)

func TestLinkedList(t *testing.T) {
	tt := []struct{
		seed []interface{}
		reversed []interface{}
	}{
		{
			seed: []interface{}{1,2,3,4,5},
			reversed: []interface{}{5,4,3,2,1},
		},
		{
			seed: []interface{}{"a","b","c","d","e"},
			reversed: []interface{}{"e","d","c","b","a"},
		},
		{
			seed: []interface{}{1},
			reversed: []interface{}{1},
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			ll := linkedlist.New()
			for _, data := range tc.seed {
				ll.Add(data)
			}

			actualSize, expectedSize := ll.Count(), len(tc.seed)
			if actualSize != expectedSize {
				t.Errorf("expected Linked List size to be %d, got %d", expectedSize, actualSize)
			}

			items := ll.Items()
			if len(items) != expectedSize {
				t.Errorf("expected Linked List items count to be %d, got %d", expectedSize, len(items))
			}

			assertSlicesEqual(t, tc.seed, items)

			ll.Reverse()
			reversedItems := ll.Items()
			if len(reversedItems) != expectedSize {
				t.Errorf("expected Linked List items count to be %d, got %d", expectedSize, len(items))
			}

			assertSlicesEqual(t, tc.reversed, reversedItems)
		})
	}
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