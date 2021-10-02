package daryheap_test

import (
	"github.com/denismitr/gds/daryheap"
	"testing"
)

func BenchmarkDaryHeap_Contains(b *testing.B) {
	dh, err := daryheap.New(2)
	if err != nil {
		b.Fatal(err)
	}

	keys := [...]string{"foo", "bar", "baz", "abc123", "cba123abc", "foobar", "barBaz", "fooBaz"}

	dh.Insert(testableString("foo"), 45.4)
	dh.Insert(testableString("bar"), 46.4)
	dh.Insert(testableString("baz"), 5.3)
	dh.Insert(testableString("abc123"), -45.4)
	dh.Insert(testableString("cba123abc"), 145.9)
	dh.Insert(testableString("foobar"), -148.9)
	dh.Insert(testableString("barBaz"), 95.932)
	dh.Insert(testableString("fooBaz"), -91.932)

	for i := 0; i < b.N; i++ {
		key := keys[i % len(keys)]
		dh.Contains(testableString(key))
	}
}
