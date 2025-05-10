package weighted

import (
	"math/rand/v2"
	"strconv"
	"testing"
)

func BenchmarkSW_Next(b *testing.B) {
	b.ReportAllocs()
	w := &SW[string]{}
	for i := 0; i < 50; i++ {
		w.Add("item-"+strconv.Itoa(i), rand.IntN(100)+100)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Next()
		}
	})
}

func BenchmarkRRW_Next(b *testing.B) {
	b.ReportAllocs()
	w := &RRW[string]{}
	for i := 0; i < 50; i++ {
		w.Add("item-"+strconv.Itoa(i), rand.IntN(100)+100)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Next()
		}
	})
}

func BenchmarkRandW_Next(b *testing.B) {
	b.ReportAllocs()
	w := NewRandW[string]()
	for i := 0; i < 50; i++ {
		w.Add("item-"+strconv.Itoa(i), rand.IntN(100)+100)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Next()
		}
	})
}
