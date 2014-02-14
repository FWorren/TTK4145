package driver

import (
	"fmt"
	"time"
)

type State_t int

const (
	RUN State_t = iota
	WAIT
	DOOR
	STOPS
	STOP_OBS
	UNDEF
	EXIT
)

var State State_t
var order_chan chan int

func Elevator_init() {
	OrderLogic_init()
	Elevator_clear_all_lights()
	if Elev_get_floor_sensor_signal() != -1 {
		State = WAIT
		OrderLogic_update_previous_order(Elev_get_floor_sensor_signal(), -1)
		Elev_set_floor_indicator(Previous_order.floor)
	}
	Elev_set_speed(-300)
	for {
		if Elev_get_floor_sensor_signal() != -1 {
			State = WAIT
			OrderLogic_update_previous_order(Elev_get_floor_sensor_signal(), -1)
			Elev_set_floor_indicator(Previous_order.floor)
			Elev_set_speed(0)
			break
		}
	}
	order_chan := make(chan int, 1)
	go OrderLogic_search_for_orders(order_chan)
}

func Elevator_statemachine() {
	for {
		switch State {
		case RUN:
			Elevator_run()
		case WAIT:
			Elevator_wait()
		case DOOR:
			Elevator_door()
		case STOPS:
			Elevator_stop()
		case STOP_OBS:
			Elevator_stop_obstruction()
		case UNDEF:
			Elevator_init()
		default:
			fmt.Println("Error! Program terminated!\n")
			State = EXIT
		}
		if State == EXIT {
			break
		}
	}
}

func Elevator_wait() {
	for {
		if OrderLogic_get_number_of_orders() > 0 {
			State = RUN
		}
		if Elev_get_stop_signal() {
			State = STOPS
		}
		if State != WAIT || State == STOPS {
			break
		}
	}
}

func Elevator_run() {
	OrderLogic_set_head_order()
	Elev_set_speed(300 * Head_order.dir)
	for {
		floor := Elev_get_floor_sensor_signal()
		if floor != -1 {
			OrderLogic_update_previous_order(floor, Head_order.dir)
			Elev_set_floor_indicator(floor)
		}
		if floor == Head_order.floor {
			State = DOOR
			Elevator_break(Head_order.dir)
			break
		}
		if Elev_get_stop_signal() {
			State = STOPS
			Elevator_break(Head_order.dir)
			break
		}
		if Elev_get_obstruction_signal() {
			State = STOP_OBS
			Elevator_break(Head_order.dir)
			break
		}
	}
}

func Elevator_door() {
	OrderLogic_delete_current_order(Head_order.floor)
	Elevator_clear_lights_current_floor(Head_order.floor)
	if Elev_get_floor_sensor_signal() != -1 {
		Elev_set_door_open_lamp(1)
	}
	time.Sleep(3 * time.Second)
	for {
		if !Elev_get_obstruction_signal() {
			break
		}
	}
	if OrderLogic_get_number_of_orders() > 0 {
		State = RUN
		Elev_set_door_open_lamp(0)
	} else {
		State = WAIT
		Elev_set_door_open_lamp(0)
	}
}

func Elevator_stop() {
	fmt.Println("The elevator has stopped!\n1. If you wish to order a new floor, do so, or.\n2. Press Ctrl + c to exit program.\n")
	Elevator_clear_all_lights()
	Elev_set_stop_lamp(1)
	OrderLogic_delete_all_orders()
	if Elev_get_floor_sensor_signal() != -1 {
		Elev_set_door_open_lamp(1)
	}
	for {
		//OrderLogic_search_for_orders()
		if OrderLogic_get_number_of_orders() > 0 {
			State = RUN
			break
		}
	}
	Elev_set_stop_lamp(0)
	Elev_set_door_open_lamp(0)
}

func Elevator_stop_obstruction() {
	for {
		if !Elev_get_obstruction_signal() {
			State = RUN
			break
		}
	}
}

func Elevator_clear_lights_current_floor(current_floor int) {
	Elev_set_button_lamp(BUTTON_COMMAND, current_floor, 0)
	if current_floor > 0 {
		Elev_set_button_lamp(BUTTON_CALL_DOWN, current_floor, 0)
	}
	if current_floor < N_FLOORS-1 {
		Elev_set_button_lamp(BUTTON_CALL_UP, current_floor, 0)
	}
}

func Elevator_clear_all_lights() {
	Elev_set_door_open_lamp(0)
	Elev_set_stop_lamp(0)
	for i := 0; i < N_FLOORS; i++ {
		Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
		if i > 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0)
		}
		if i < N_FLOORS-1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, i, 0)
		}
	}
}

func Elevator_break(direction int) {
	Elev_set_speed(40 * (-direction))
	time.Sleep(2 * time.Millisecond)
	Elev_set_speed(0)
}
