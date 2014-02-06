package main

import (
	"fmt"
	"./network"
	"./driver"
)

func main() {
	 // Initialize hardware
    if !Elev_init() {
        fmt.Println("Unable to initialize elevator hardware\n");
        return 1;
    }
    fmt.Println("Press STOP button to stop elevator and exit program.\n");

    Elev_set_speed(300)

    for {
    	if Elev_get_floor_sensor_signal() {
    		Elev_set_speed(0)
    	}
    }
}