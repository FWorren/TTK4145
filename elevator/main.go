package main

import (
	driver "./driver"
	//network "./network"
	"fmt"
)

func main() {
	// Initialize hardware
	if driver.Elev_init() == 0 {
		fmt.Println("Unable to initialize elevator hardware\n")
	}

	fmt.Println("Press STOP button to stop elevator and exit program.\n")

	driver.Elevator_statemachine()
	
	//network.Network()
}