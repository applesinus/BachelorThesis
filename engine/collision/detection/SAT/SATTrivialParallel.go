package sat

import (
	"BachelorThesis/engine/objects"
	"log"
)

func SATTrivialParallel(aID, bID int, objectPool *[]objects.Object, resolveAlgorithm string) {
	switch (*objectPool)[aID].(type) {
	case *objects.Sphere:
		switch (*objectPool)[bID].(type) {
		case *objects.Sphere:
			satSphereSphere(aID, bID, objectPool, resolveAlgorithm)
		}

	default:
		log.Panicf("Unknown object type: %T", (*objectPool)[aID])
	}
}
