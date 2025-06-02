package main

import (
	"BachelorThesis/engine"
	"BachelorThesis/engine/constants"
	"context"
	"fmt"
	"log"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var algorithm, secondaryAlgorithm, resolveAlgorithm, pipeline string

	ctx, cancel := context.WithCancel(context.Background())

	command := "start"
	for command != "exit" {
		if command == "start" {
			fmt.Printf("\n ===== SIMULATION OPTIONS  =====\n")

			fmt.Printf("Avaliable algorithms:\n")
			fmt.Printf("\t1. %s%s\n", constants.SaP, constants.N)
			fmt.Printf("\t2. %s%s\n", constants.SaP, constants.PT)
			fmt.Printf("\t3. %s%s\n", constants.SaP, constants.PNT)
			fmt.Printf("Enter a number to choose an algorithm (1/2/3): ")
			algo := 0
			for algo == 0 {
				_, err := fmt.Scanln(&algo)
				if err != nil || algo < 1 || algo > 3 {
					algo = 0
					continue
				}
			}
			switch algo {
			case 1:
				algorithm = constants.SaP
				constants.AlgoType = constants.N
			case 2:
				algorithm = constants.SaP
				constants.AlgoType = constants.PT
			case 3:
				algorithm = constants.SaP
				constants.AlgoType = constants.PNT
			default:
				log.Panicf("Unknown algorithm: %d", algo)
				continue
			}

			fmt.Printf("Avaliable secondary algorithms:\n")
			fmt.Printf("\t1. %s%s\n", constants.SAT, constants.N)
			fmt.Printf("\t2. %s%s\n", constants.SAT, constants.PT)
			fmt.Printf("Enter a number to choose an algorithm (1/2): ")
			secAlgo := 0
			for secAlgo == 0 {
				_, err := fmt.Scanln(&secAlgo)
				if err != nil || secAlgo < 1 || secAlgo > 2 {
					secAlgo = 0
					continue
				}
			}
			switch secAlgo {
			case 1:
				secondaryAlgorithm = constants.SAT
				constants.SecondaryAlgoType = constants.N
			case 2:
				secondaryAlgorithm = constants.SAT
				constants.SecondaryAlgoType = constants.PT
			default:
				log.Panicf("Unknown algorithm: %d", secAlgo)
				continue
			}

			fmt.Printf("Avaliable resolve algorithms:\n")
			fmt.Printf("\t1. %s%s\n", constants.PGS, constants.N)
			fmt.Printf("\t2. %s%s\n", constants.PGS, constants.PNT)
			fmt.Printf("Enter a number to choose an algorithm (1/2): ")
			resAlgo := 0
			for resAlgo == 0 {
				_, err := fmt.Scanln(&resAlgo)
				if err != nil || resAlgo < 1 || resAlgo > 2 {
					resAlgo = 0
					continue
				}
			}
			switch resAlgo {
			case 1:
				resolveAlgorithm = constants.PGS
				constants.ResolveAlgoType = constants.N
			case 2:
				resolveAlgorithm = constants.PGS
				constants.ResolveAlgoType = constants.PNT
			default:
				log.Panicf("Unknown algorithm: %d", resAlgo)
				continue
			}

			pipeline = ""
			for pipeline != "y" && pipeline != "n" {
				fmt.Printf("Would you like to use parallel pipeline? (y/n): ")
				fmt.Scanln(&pipeline)
			}
			if pipeline == "y" {
				constants.Pipeline = constants.ParallelPipeline
			} else {
				constants.Pipeline = constants.SequentialPipeline
			}

			ctx, cancel = context.WithCancel(context.Background())
			go engine.Run(algorithm, secondaryAlgorithm, resolveAlgorithm, ctx, cancel)
		}

		fmt.Printf("\n ===== ENTER A COMMAND  =====\n")

		fmt.Scanln(&command)

		if command == "end" || command == "exit" {
			select {
			case <-ctx.Done():
				log.Printf("Simulation already ended")
				continue
			default:
				cancel()
			}
		}
	}

	log.Printf("Simulation ended")
}
