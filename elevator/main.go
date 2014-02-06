package main

import (
	"fmt"
	network "./network"
	driver "./driver"
)

func main() {
	 // Initialize hardware
    if driver.Elev_init() == 0 {
        fmt.Println("Unable to initialize elevator hardware\n");
    }
    fmt.Println("Press STOP button to stop elevator and exit program.\n");

    network.Network()

    driver.Elev_set_speed(300)

    for {
    	if driver.Elev_get_floor_sensor_signal() == 1 {
    		driver.Elev_set_speed(0)
    	}
    }
}