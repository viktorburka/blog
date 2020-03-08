package bubblesort

func BubbleSort(data []int) []int {
	size := len(data)
	for i := 0; i < size; i++ {
		for j := 0; j < size-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
	return data
}
