package engine

import (
	"BachelorThesis/engine/objects"
	"BachelorThesis/engine/visualizer"
	"log"
)

func Run(algorithm string, endChan chan bool) {
	log.Printf("Simulation of %s started", algorithm)

	pool := make([]objects.Object, 0)

	visualizer.Start(algorithm, endChan, &pool)

	log.Printf("Simulation of %s ended", algorithm)
}
