package main

import (
	"fmt"
	"./network"
	"./driver"
)

func main() {
	 // Initialize hardware
    if !elev_init() {
        fmt.Println("Unable to initialize elevator hardware\n");
        return 1;
    }
    fmt.Println("Press STOP button to stop elevator and exit program.\n");
}