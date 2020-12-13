package mapper

import (
	"math/rand"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	rand.Seed(1)
	var data = New(10)
	for i := 0; i < b.N; i++ {
		data.Insert([]string{"foo.bar"}, 1)
	}
}

func BenchmarkRemove(b *testing.B) {
	rand.Seed(1)
	var data = New(10)
	for i := 0; i < b.N; i++ {
		data.Remove("foo.bar")
	}
}

func BenchmarkGetResults(b *testing.B) {
	rand.Seed(1)
	var data = New(10)
	for i := 0; i < b.N; i++ {
		data.GetResults()
	}
}
