package SaP

import (
	sat "BachelorThesis/engine/collision/detection/SAT"
	tgs "BachelorThesis/engine/collision/resolving/TGS"
	"BachelorThesis/engine/constants"
	"BachelorThesis/engine/objects"
	"log"
	"math"
	"runtime"
	"sync"
)

type intPair struct {
	a int
	b int
}

func Collision(objectPool *[]objects.Object, secondaryAlgorithm, resolveAlgorithm string) {
	// First step: sort
	if constants.AlgoType == constants.N {
		quickSort(objectPool)
	} else {
		radixSort(objectPool)
	}

	if !isSorted(*objectPool) {
		log.Panicf("Objects are not sorted OLOLO")
	}

	// Second step: Sweep and Prune
	if constants.AlgoType == constants.PNT {
		sapParallelNonTrivial(objectPool, secondaryAlgorithm, resolveAlgorithm)
	} else {
		sapNoParallel(objectPool, secondaryAlgorithm, resolveAlgorithm)
	}
}

// --- No parallel algorithm ---

func sapNoParallel(objectPool *[]objects.Object, secondaryAlgorithm, resolveAlgorithm string) {
	if len(*objectPool) <= 1 {
		return
	}

	activeObjects := make(map[int]*objects.Object, 0)
	pairs := make([]intPair, 0)
	toDel := make([]int, len(*objectPool))

	for a, obj := range *objectPool {
		bb, err := obj.GetBoundingBox()
		if err != nil {
			log.Printf("Warning: failed to get bounding box for object %s: %v", obj.GetId(), err)
			continue
		}

		toDel = toDel[:0]

		for b, activeObj := range activeObjects {
			activeBB, err := (*activeObj).GetBoundingBox()
			if err != nil {
				log.Printf("Warning: failed to get bounding box for object %s: %v", (*activeObj).GetId(), err)
				continue
			}

			if bb.Min.X < activeBB.Max.X {
				if checkOverlapYZ(&obj, activeObj) {
					if constants.Pipeline == constants.ParallelPipeline {
						switch secondaryAlgorithm {
						case constants.SAT:
							switch constants.SecondaryAlgoType {
							case constants.N:
								sat.SATNoParallel(a, b, objectPool, resolveAlgorithm)
							case constants.PT:
								sat.SATTrivialParallel(a, b, objectPool, resolveAlgorithm)
							default:
								log.Panicf("Unknown secondary algorithm type: %s", constants.SecondaryAlgoType)
							}

						// if there is no secondary algorithm
						case constants.NoAlgo:
							switch resolveAlgorithm {
							case constants.PGS:
								switch constants.ResolveAlgoType {
								case constants.N:
									tgs.TGSNoParallel(a, b, objectPool)
								case constants.PNT:
									// TODO
									//pgs.PGSParallelNonTrivial(a, b, objectPool)
								default:
									log.Panicf("Unknown resolve algorithm type: %s", constants.ResolveAlgoType)
								}

							// if there is no resolve algorithm just return
							case constants.NoAlgo:
								return
							default:
								log.Panicf("Unknown resolve algorithm: %s", resolveAlgorithm)
							}
						default:
							log.Panicf("Unknown secondary algorithm: %s", secondaryAlgorithm)
						}
					} else {
						pairs = append(pairs, intPair{a: a, b: b})
					}
				}
			} else {
				toDel = append(toDel, b)
			}
		}

		for _, d := range toDel {
			delete(activeObjects, d)
		}

		activeObjects[a] = &(*objectPool)[a]
	}

	if constants.Pipeline == constants.ParallelPipeline {
		return
	}

	switch secondaryAlgorithm {
	case constants.SAT:
		switch constants.SecondaryAlgoType {
		case constants.N:
			for _, pair := range pairs {
				sat.SATNoParallel(pair.a, pair.b, objectPool, resolveAlgorithm)
			}
		case constants.PT:
			for _, pair := range pairs {
				sat.SATTrivialParallel(pair.a, pair.b, objectPool, resolveAlgorithm)
			}
		default:
			log.Panicf("Unknown secondary algorithm type: %s", constants.SecondaryAlgoType)
		}

	// if there is no secondary algorithm
	case constants.NoAlgo:
		switch resolveAlgorithm {
		case constants.PGS:
			switch constants.ResolveAlgoType {
			case constants.N:
				for _, pair := range pairs {
					tgs.TGSNoParallel(pair.a, pair.b, objectPool)
				}
			case constants.PNT:
				// TODO
				//pgs.PGSParallelNonTrivial(a, b, objectPool)
			default:
				log.Panicf("Unknown resolve algorithm type: %s", constants.ResolveAlgoType)
			}

		// if there is no resolve algorithm just return
		case constants.NoAlgo:
			return
		default:
			log.Panicf("Unknown resolve algorithm: %s", resolveAlgorithm)
		}
	default:
		log.Panicf("Unknown secondary algorithm: %s", secondaryAlgorithm)
	}
}

// --- Parallel Non Trivial algorithm ---

func sapParallelNonTrivial(objectPool *[]objects.Object, secondaryAlgorithm, resolveAlgorithm string) {
	if len(*objectPool) <= 1 {
		return
	}

	workersCount := runtime.NumCPU() - 1
	if workersCount < 3 {
		log.Printf("Warning: number of workers is less than 3: %d. Using no parallel algorithm", workersCount)
		sapNoParallel(objectPool, secondaryAlgorithm, resolveAlgorithm)
	}
	wg := new(sync.WaitGroup)
	wg.Add(workersCount)

	outChan := make(chan *intPair, int(math.Sqrt(float64(len(*objectPool)))))

	go func() {
		wg.Wait()
		close(outChan)
	}()

	for i := 0; i < workersCount; i++ {
		start := i * len(*objectPool) / workersCount
		end := (i + 1) * len(*objectPool) / workersCount
		if end > len(*objectPool) {
			end = len(*objectPool)
		}

		go func(start, end int) {
			defer wg.Done()

			for i, obj := range (*objectPool)[start:end] {
				bb, err := obj.GetBoundingBox()
				if err != nil {
					log.Printf("Warning: failed to get bounding box for object %s: %v", obj.GetId(), err)
					continue
				}

				for j, activeObj := range (*objectPool)[start+i+1:] {
					activeBB, err := activeObj.GetBoundingBox()
					if err != nil {
						log.Printf("Warning: failed to get bounding box for object %s: %v", activeObj.GetId(), err)
						continue
					}

					if bb.Min.X < activeBB.Max.X {
						if checkOverlapYZ(&obj, &activeObj) {
							if constants.Pipeline == constants.ParallelPipeline {
								switch secondaryAlgorithm {
								case constants.SAT:
									switch constants.SecondaryAlgoType {
									case constants.N:
										sat.SATNoParallel(start+i, start+i+j+1, objectPool, resolveAlgorithm)
									case constants.PT:
										sat.SATTrivialParallel(start+i, start+i+j+1, objectPool, resolveAlgorithm)
									default:
										log.Panicf("Unknown secondary algorithm type: %s", constants.SecondaryAlgoType)
									}
								// if there is no secondary algorithm
								case constants.NoAlgo:
									switch resolveAlgorithm {
									case constants.PGS:
										switch constants.ResolveAlgoType {
										case constants.N:
											tgs.TGSNoParallel(start+i, start+j+1, objectPool)
										case constants.PNT:
											// TODO
											//pgs.PGSParallelNonTrivial(start+i, start+j+1, objectPool)
										default:
											log.Panicf("Unknown resolve algorithm type: %s", constants.ResolveAlgoType)
										}

									// if there is no resolve algorithm just return
									case constants.NoAlgo:
										return
									default:
										log.Panicf("Unknown resolve algorithm: %s", resolveAlgorithm)
									}
								}
							} else {
								outChan <- &intPair{a: start + i, b: start + j + i + 1}
							}
						}
					} else {
						break
					}
				}
			}
		}(start, end)
	}

	pairs := make([]intPair, 0)
	for pair := range outChan {
		pairs = append(pairs, *pair)
	}

	if constants.Pipeline == constants.ParallelPipeline {
		return
	}

	switch secondaryAlgorithm {
	case constants.SAT:
		switch constants.SecondaryAlgoType {
		case constants.N:
			for _, pair := range pairs {
				sat.SATNoParallel(pair.a, pair.b, objectPool, resolveAlgorithm)
			}
		case constants.PT:
			for _, pair := range pairs {
				sat.SATTrivialParallel(pair.a, pair.b, objectPool, resolveAlgorithm)
			}
		default:
			log.Panicf("Unknown secondary algorithm type: %s", constants.SecondaryAlgoType)
		}
	// if there is no secondary algorithm
	case constants.NoAlgo:
		switch resolveAlgorithm {
		case constants.PGS:
			switch constants.ResolveAlgoType {
			case constants.N:
				for _, pair := range pairs {
					tgs.TGSNoParallel(pair.a, pair.b, objectPool)
				}
			case constants.PNT:
				// TODO
				//pgs.PGSParallelNonTrivial(start+i, start+j+1, objectPool)
			default:
				log.Panicf("Unknown resolve algorithm type: %s", constants.ResolveAlgoType)
			}

		// if there is no resolve algorithm just return
		case constants.NoAlgo:
			return
		default:
			log.Panicf("Unknown resolve algorithm: %s", resolveAlgorithm)
		}
	}
}

// --- Helper function ---

func checkOverlapYZ(objA, objB *objects.Object) bool {
	bbA, errA := (*objA).GetBoundingBox()
	bbB, errB := (*objB).GetBoundingBox()

	if errA != nil || errB != nil {
		return false
	}

	yOverlap := (bbA.Max.Y >= bbB.Min.Y && bbA.Min.Y <= bbB.Max.Y) || (bbB.Max.Y >= bbA.Min.Y && bbB.Min.Y <= bbA.Max.Y)
	zOverlap := (bbA.Max.Z >= bbB.Min.Z && bbA.Min.Z <= bbB.Max.Z) || (bbB.Max.Z >= bbA.Min.Z && bbB.Min.Z <= bbA.Max.Z)

	return yOverlap && zOverlap
}
