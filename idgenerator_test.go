package axiom

import (
	"testing"
	"fmt"
	"runtime"
)

func TestIdGen_Next(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	workerID := []int64{1}
	for _, wid := range workerID {
		idgen := NewID(wid)
		for i := 1; i <= 10; i++ {
			fmt.Printf("%d: %d \n", wid, idgen.Next())
		}
	}
}

func BenchmarkIdGen_Next(b *testing.B) {
	idgen := NewID(1)
	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		idgen.Next()
	}
}
