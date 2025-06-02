package SaP

import (
	"BachelorThesis/engine/objects"
)

func quickSort(objects *[]objects.Object) {
	if len(*objects) <= 1 {
		return
	}
	quickSortRecursive(*objects, 0, len(*objects)-1)
}

func quickSortRecursive(objects []objects.Object, low, high int) {
	if low < high {
		pivotIndex := partition(objects, low, high)

		quickSortRecursive(objects, low, pivotIndex-1)
		quickSortRecursive(objects, pivotIndex+1, high)
	}
}

func partition(objects []objects.Object, low, high int) int {
	pivotBB, _ := objects[high].GetBoundingBox()
	pivotValue := pivotBB.Min.X

	i := low - 1

	for j := low; j < high; j++ {
		currentBB, _ := objects[j].GetBoundingBox()
		if currentBB.Min.X <= pivotValue {
			i++
			objects[i], objects[j] = objects[j], objects[i]
		}
	}

	objects[i+1], objects[high] = objects[high], objects[i+1]
	return i + 1
}
