package collision

import (
	"BachelorThesis/engine/collision/detection/SaP"
	"BachelorThesis/engine/constants"
	"BachelorThesis/engine/objects"
	"log"
)

func ProcessCollisions(objects *[]objects.Object, algorithm, secondaryAlgorithm, resolveAlgorithm string) {
	switch algorithm {
	case constants.SaP:
		SaP.Collision(objects, secondaryAlgorithm, resolveAlgorithm)
	case constants.NoAlgo:
		return
	default:
		log.Panicf("Unknown algorithm: %s", algorithm)
	}
}
