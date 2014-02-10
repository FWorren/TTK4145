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

    driver.State = driver.UNDEF

    driver.Elevator_statemachine()

    /*driver.Elev_set_speed(300)

    for {
        floor := driver.Elev_get_floor_sensor_signal() 
    	if floor != -1 {
            driver.Elev_set_floor_indicator(floor)
    		driver.Elev_set_speed(-300)
    	}
        if driver.Elev_get_stop_signal() {
            driver.Elev_set_speed(0)
            break
        }
    } */
}