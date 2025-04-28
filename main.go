package main

import (
	"BachelorThesis/engine"
	"BachelorThesis/engine/constants"
	"fmt"
)

func main() {
	endChan := make(chan bool)

	go engine.Run(constants.NoAlgo, endChan)

	command := ""
	for command != "end" && command != "exit" {
		fmt.Scanln(&command)
		if command == "end" || command == "exit" {
			endChan <- true
		}

		if command == "start" {
			go engine.Run(constants.NoAlgo, endChan)
		}
	}
}
