package main

import "fmt"

func main() {
	comp := func(a, b float64) bool {
		return a < b
	}
	var _ Less[float64] = comp

	h := NewMinHeap(comp)
	arr := []float64{3, 1, 6, 5, 2, 4, 1}
	for _, v := range arr {
		h.Insert(v)
	}

	fmt.Printf("Heap: %s\n", h)
	fmt.Println()

	for h.Len() > 0 {
		fmt.Printf("Min: %v\n", h.ExtractMin())
		fmt.Printf("Heap: %s\n", h)
	}
}
