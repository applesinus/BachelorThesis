package SaP

import (
	"BachelorThesis/engine/objects"
	"math"
	"runtime"
	"sync"
)

const (
	RADIX_BITS = 8
	RADIX_SIZE = 1 << RADIX_BITS
	RADIX_MASK = RADIX_SIZE - 1
)

type sortItem struct {
	obj objects.Object
	key uint64
}

func floatToSortableUint64(f float64) uint64 {
	bits := math.Float64bits(f)
	if bits&(1<<63) != 0 {
		return ^bits
	} else {
		return bits | (1 << 63)
	}
}

func countingSort(items []sortItem, shift uint, temp []sortItem) {
	count := make([]int, RADIX_SIZE)
	for i := range items {
		radix := (items[i].key >> shift) & RADIX_MASK
		count[radix]++
	}

	for i := 1; i < RADIX_SIZE; i++ {
		count[i] += count[i-1]
	}

	for i := len(items) - 1; i >= 0; i-- {
		radix := (items[i].key >> shift) & RADIX_MASK
		count[radix]--
		temp[count[radix]] = items[i]
	}

	copy(items, temp)
}

func parallelCountingSort(items []sortItem, shift uint, temp []sortItem, numWorkers int) {
	n := len(items)
	if n <= 10000 || numWorkers <= 1 {
		countingSort(items, shift, temp)
		return
	}

	chunkSize := (n + numWorkers - 1) / numWorkers

	localCounts := make([][]int, numWorkers)
	for i := range localCounts {
		localCounts[i] = make([]int, RADIX_SIZE)
	}

	var wg sync.WaitGroup

	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			start := w * chunkSize
			end := start + chunkSize
			if end > n {
				end = n
			}

			for i := start; i < end; i++ {
				radix := (items[i].key >> shift) & RADIX_MASK
				localCounts[w][radix]++
			}
		}(worker)
	}
	wg.Wait()

	globalCount := make([]int, RADIX_SIZE)
	for i := 0; i < RADIX_SIZE; i++ {
		for w := 0; w < numWorkers; w++ {
			globalCount[i] += localCounts[w][i]
		}
	}

	for i := 1; i < RADIX_SIZE; i++ {
		globalCount[i] += globalCount[i-1]
	}

	for i := n - 1; i >= 0; i-- {
		radix := (items[i].key >> shift) & RADIX_MASK
		globalCount[radix]--
		temp[globalCount[radix]] = items[i]
	}

	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			start := w * chunkSize
			end := start + chunkSize
			if end > n {
				end = n
			}
			copy(items[start:end], temp[start:end])
		}(worker)
	}
	wg.Wait()
}

func radixSort(objects *[]objects.Object) {
	if len(*objects) <= 1 {
		return
	}

	numWorkers := runtime.NumCPU()
	n := len(*objects)

	items := make([]sortItem, n)
	temp := make([]sortItem, n)

	chunkSize := (n + numWorkers - 1) / numWorkers
	var wg sync.WaitGroup

	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			start := w * chunkSize
			end := start + chunkSize
			if end > n {
				end = n
			}

			for i := start; i < end; i++ {
				bb, _ := (*objects)[i].GetBoundingBox()
				items[i] = sortItem{
					obj: (*objects)[i],
					key: floatToSortableUint64(bb.Min.X),
				}
			}
		}(worker)
	}
	wg.Wait()

	for shift := uint(0); shift < 64; shift += RADIX_BITS {
		parallelCountingSort(items, shift, temp, numWorkers)
	}

	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			start := w * chunkSize
			end := start + chunkSize
			if end > n {
				end = n
			}

			for i := start; i < end; i++ {
				(*objects)[i] = items[i].obj
			}
		}(worker)
	}
	wg.Wait()
}

func isSorted(objects []objects.Object) bool {
	for i := 1; i < len(objects); i++ {
		bb1, _ := objects[i-1].GetBoundingBox()
		bb2, _ := objects[i].GetBoundingBox()
		if bb1.Min.X > bb2.Min.X {
			return false
		}
	}
	return true
}
