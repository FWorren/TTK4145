package driver

import (
	"fmt"
	"time"
)

type State_t int

const (
	RUN State_t = iota
	DOOR
	STOPS
	STOP_OBS
	UNDEF
)

func Elevator_eventHandler(head_order_c chan Order, prev_order_c chan Order) {
	//var State State_t
	//State = UNDEF
	floor_reached := make(chan bool)
	obstruction := make(chan bool)
	stop := make(chan bool)
	Elevator_init()
	for {
		select {
		case <-floor_reached:
			//State = DOOR
			Elevator_door()
		case new_order := <-head_order_c:
			//State = RUN
			Elevator_run(floor_reached, new_order, obstruction, stop, prev_order_c)
		case <-obstruction:
			//State = STOP_OBS
			Elevator_stop_obstruction()
		case <-stop:
			//State = STOPS
			Elevator_stop()
		}
	}
}

/*func Elevator_statemachine(state State_t, floor_reached chan bool, head_order_c chan Order, obstruction chan bool, stop chan bool, prev_order_c chan Order) {
	switch state {
	case RUN:
		Elevator_run(floor_reached, head_order_c, obstruction, stop, prev_order_c)
	case DOOR:
		Elevator_door()
	case STOPS:
		Elevator_stop()
	case STOP_OBS:
		Elevator_stop_obstruction()
	case UNDEF:
		prev_order_c <- Elevator_init(floor_reached)
	}
}*/

func Elevator_init() {
	Elevator_clear_all_lights()
	Elev_set_speed(-300)
	for {
		time.Sleep(10 * time.Millisecond)
		floor := Elev_get_floor_sensor_signal()
		fmt.Println(floor)
		if floor != -1 {
			Elevator_break(-1)
			/*var current Order
			current.Floor = floor
			current.Dir = -1*/
			break
		}
	}
	fmt.Println("init complete ")
}

func Elevator_run(floor_reached chan bool, head_order Order, obstruction chan bool, stop chan bool, prev_order chan Order) {
	Elev_set_speed(300 * head_order.Dir)
	for {
		time.Sleep(25 * time.Millisecond)
		current_floor := Elev_get_floor_sensor_signal()
		if current_floor != -1 {
			var current Order
			current.Floor = current_floor
			current.Dir = head_order.Dir
			prev_order <- current
			Elev_set_floor_indicator(current_floor)
		}
		if current_floor == head_order.Floor {
			Elevator_break(head_order.Dir)
			floor_reached <- true
			break
		}
		if Elev_get_stop_signal() {
			Elevator_break(head_order.Dir)
			stop <- true
			break
		}
		if Elev_get_obstruction_signal() {
			Elevator_break(head_order.Dir)
			obstruction <- true
			break
		}
	}
}

func Elevator_door() {
	if Elev_get_floor_sensor_signal() != -1 {
		Elev_set_door_open_lamp(1)
		time.Sleep(3 * time.Second)
		for {
			time.Sleep(25 * time.Millisecond)
			if !Elev_get_obstruction_signal() {
				return
			}
		}
		Elev_set_door_open_lamp(0)
	} else {

	}
}

func Elevator_stop() {
	fmt.Println("The elevator has stopped!\n1. If you wish to order a new floor, do so, or.\n2. Press Ctrl + c to exit program.\n")
	Elevator_clear_all_lights()
	Elev_set_stop_lamp(1)
	for {
		time.Sleep(25 * time.Millisecond)
	}
}

func Elevator_stop_obstruction() {
	for {
		time.Sleep(25 * time.Millisecond)
		if !Elev_get_obstruction_signal() {
			return
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
	Elev_set_speed(20 * (-direction))
	//time.Sleep(10 * time.Millisecond)
	Elev_set_speed(0)
}
