package bubblesort

import (
	"math/rand"
	"testing"
)

func randomData(size int) []int {
	d := make([]int, size, size)
	for i := 0; i < size; i++ {
		d[i] = rand.Intn(size)
	}
	return d
}

func BenchmarkBubbleSort(b *testing.B) {
	data := randomData(10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
        BubbleSort(data)
    }
}
