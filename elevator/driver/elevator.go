package driver

import (
	"fmt"
	"time"
)

type State_t int

const (
	RUN State_t = iota
	DOOR
	WAIT
	STOPS
	STOP_OBS
	UNDEF
)

func Elevator_eventHandler(head_order_c chan Order, prev_order_c chan Order, del_order chan Order, state_c chan State_t) {
	floor_reached := make(chan bool)
	obstruction := make(chan bool)
	stop := make(chan bool)
	get_prev_floor := make(chan Order)
	delete_order := make(chan Order)
	state := make(chan State_t)
	var head_order Order
	state_c <- WAIT

	for {
		time.Sleep(10 * time.Millisecond)
		select {
		case head_order = <-head_order_c:
			go Elevator_run(floor_reached, head_order, obstruction, stop, prev_order_c, state)
		case <-floor_reached:
			go Elevator_door(head_order, delete_order, state)
		case <-obstruction:
			go Elevator_stop_obstruction()
		case <-stop:
			go Elevator_stop()
		case update_prev := <-get_prev_floor:
			prev_order_c <- update_prev
		case del_req := <-delete_order:
			del_order <- del_req
		case new_state := <-state:
			state_c <- new_state
		}
	}
}

func Elevator_init() (init bool, prev Order) {
	Elevator_clear_all_lights()
	Elev_set_speed(-300)
	for {
		time.Sleep(10 * time.Millisecond)
		floor := Elev_get_floor_sensor_signal()
		if floor != -1 {
			Elev_set_floor_indicator(floor)
			Elevator_break(-1)
			set := true
			var current Order
			current.Floor = floor
			current.Dir = -1
			return set, current
		}
	}
}

func Elevator_run(floor_reached chan bool, head_order Order, obstruction chan bool, stop chan bool, get_prev_floor chan Order, state chan State_t) {
	Elev_set_speed(300 * head_order.Dir)
	state <- RUN
	for {
		time.Sleep(10 * time.Millisecond)
		current_floor := Elev_get_floor_sensor_signal()
		if current_floor != -1 {
			var current Order
			current.Floor = current_floor
			current.Dir = head_order.Dir
			get_prev_floor <- current
			Elev_set_floor_indicator(current_floor)
		}
		if current_floor == head_order.Floor {
			fmt.Println("FLOOR REEEEEACHED")
			Elevator_break(head_order.Dir)
			floor_reached <- true
			return
		}
		if Elev_get_stop_signal() {
			Elevator_break(head_order.Dir)
			stop <- true
			return
		}
		if Elev_get_obstruction_signal() {
			Elevator_break(head_order.Dir)
			obstruction <- true
			return
		}
	}
}

func Elevator_door(head_order Order, delete_order chan Order, state chan State_t) {
	if Elev_get_floor_sensor_signal() != -1 {
		Elev_set_door_open_lamp(1)
		state <- DOOR
		time.Sleep(3 * time.Second)
		for {
			time.Sleep(25 * time.Millisecond)
			if !Elev_get_obstruction_signal() {
				break
			}
		}
		state <- WAIT
		Elev_set_door_open_lamp(0)
		Elevator_clear_lights_current_floor(head_order.Floor)
		delete_order <- head_order
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
	Elev_set_speed(100 * (-direction))
	time.Sleep(20 * time.Millisecond)
	Elev_set_speed(0)
}
