package bubblesort

import (
	"math/rand"
	"strconv"
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
	tests := []int{ 10000, 20000, 30000, 40000, 50000 }
	for _, test := range tests {
		b.Run(strconv.Itoa(test), func(pb *testing.B) {
			data := randomData(test)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				BubbleSort(data)
			}
		})
	}
}
