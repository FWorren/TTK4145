package controll

import (
	"fmt"
)

type State_t int

const  (
    RUN State_t = iota
    WAIT
    DOOR
    STOP
    UNDEF
    EXIT
)

var state State_t

func Elevator_init() {
	OrderLogic_init()
	Elevator_clear_all_lights()
	if driver.Elev_get_floor_sensor_signal() != -1 {
		state := WAIT
		OrderLogic_update_previous_order(driver.Elev_get_floor_sensor_signal(),-1)
		driver.Elev_set_floor_indicator(previous_order.floor)
	}
}

func Elevator_statemachine() {
	for {
		switch state {
			case RUN:
				Elevator_run()
			case WAIT:
				Elevator_wait()
			case DOOR:
				Elevator_door()
			case STOP:
				Elevator_stop()
			case UNDEF:
				Elevator_init()
			default:
				fmt.Println("Error! Program terminated!\n")
				state = EXIT
		}
		if state == EXIT {break}
	}
}

func Elevator_wait() {
	for {
		OrderLogic_search_for_orders()
		if OrderLogic_get_number_of_internal_orders() > 0 {state = RUN}
		if state != WAIT{break}
	}
}

func Elevator_run() {
	for {
		if state != RUN{break}
	}
}

func Elevator_door() {
	for {
		if state != DOOR{break}
	}
}

func Elevator_stop() {
	for {
		if state != STOP{break}
	}
}

func Elevator_clear_lights_current_floor(current_floor int) {
	
}

func Elevator_clear_all_lights() {
	
}

func Elevator_break(direction int) {
	driver.Elev_set_speed(20*(-direction))
}

func Elevator_timer(sec float64) {
	
}