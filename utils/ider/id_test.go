package ider

import (
	"fmt"
	"runtime"
	"testing"
)

func TestID(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	id := ID()
	for i := 1; i <= 1000000; i++ {
		fmt.Printf("ID: %s \n", <-id)
	}

}

func TestIDGen_Next(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	workerID := []int64{1}
	for _, wid := range workerID {
		idgen := NewID(wid)
		for i := 1; i <= 10; i++ {
			fmt.Printf("%d: %d \n", wid, idgen.Next())
		}
	}
}

func BenchmarkID(b *testing.B) {
	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ID()
	}
}

func BenchmarkIDGen_Next(b *testing.B) {
	idgen := NewID(1)
	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		idgen.Next()
	}
}
