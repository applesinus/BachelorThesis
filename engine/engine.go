package engine

import (
	"BachelorThesis/engine/visualizer"
	"log"
)

func Run(algorithm string, endChan chan bool) {
	log.Printf("Simulation of %s started", algorithm)

	visualizer.Start(algorithm, endChan)

	log.Printf("Simulation of %s ended", algorithm)
}
