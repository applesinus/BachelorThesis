package engine

import (
	"BachelorThesis/engine/constants"
	"BachelorThesis/engine/objects"
	st "BachelorThesis/engine/singletone"
	"BachelorThesis/engine/visualizer"
	"context"
	"log"
)

var singletone *st.Engine

func Run(algorithm, secondaryAlgorithm, resolveAlgorithm string, ctx context.Context, cancel context.CancelFunc) {
	log.Printf("Simulation of %s%s + %s%s + %s%s started", algorithm, constants.AlgoType, secondaryAlgorithm, constants.SecondaryAlgoType, resolveAlgorithm, constants.ResolveAlgoType)

	if algorithm == constants.NoAlgo {
		log.Printf("Secondary algorithm is not specified, must be an error")
	}
	if secondaryAlgorithm == constants.NoAlgo {
		log.Printf("Secondary algorithm is not specified, the simulation will be speculative")
	}
	if resolveAlgorithm == constants.NoAlgo {
		log.Printf("Resolve algorithm is not specified, must be an error")
	}

	pool := make([]objects.Object, 0)

	singletone = st.NewEngine(algorithm, secondaryAlgorithm, resolveAlgorithm, &pool, ctx)
	go singletone.StartEngineLoop()

	visualizer.Start(singletone, cancel)

	log.Printf("Simulation of %s%s + %s%s + %s%s ended", algorithm, constants.AlgoType, secondaryAlgorithm, constants.SecondaryAlgoType, resolveAlgorithm, constants.ResolveAlgoType)
}
