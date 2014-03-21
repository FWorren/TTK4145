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

func Elevator_eventHandler(head_order_c chan Order, prev_order_c chan Order, del_order chan Order, state_c chan State_t, local_list [3][4]bool) {
	floor_reached := make(chan Order, 1)
	obstruction := make(chan bool, 1)
	stop := make(chan bool, 1)
	get_prev_floor_c := make(chan Order, 1)
	delete_order := make(chan Order, 1)
	state := make(chan State_t, 1)
	var prev_order Order
	prev_order.Floor = -1
	var head_order Order
	state_c <- WAIT

	for {
		time.Sleep(10 * time.Millisecond)
		select {
		case head_order = <-head_order_c:
			if head_order.Floor == prev_order.Floor {
				go Elevator_door(head_order, delete_order, state)
			} else {
				go Elevator_run(floor_reached, head_order, obstruction, stop, get_prev_floor_c, state, local_list)
			}
		case reached := <-floor_reached:
			go Elevator_door(reached, delete_order, state)
		case <-obstruction:
			go Elevator_stop_obstruction(head_order_c, head_order, state)
		case <-stop:
			go Elevator_stop(state)
		case prev_order = <-get_prev_floor_c:
			prev_order_c <- prev_order
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

func Elevator_run(floor_reached chan Order, head_order Order, obstruction chan bool, stop chan bool, get_prev_floor_c chan Order, state chan State_t, local_list [3][4]bool) {
	Elev_set_speed(300 * head_order.Dir)
	state <- RUN
	for {
		time.Sleep(10 * time.Millisecond)
		current_floor := Elev_get_floor_sensor_signal()
		if current_floor != -1 {
			var current Order
			current.Floor = current_floor
			current.Dir = head_order.Dir
			get_prev_floor_c <- current
			Elev_set_floor_indicator(current_floor)
			flag, convenient := OrderHandler_check_convenient_order(local_list, current)
			if flag {
				Elevator_break(head_order.Dir)
				floor_reached <- convenient
				return
			}
		}
		if current_floor == head_order.Floor {
			Elevator_break(head_order.Dir)
			floor_reached <- head_order
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
		Elev_set_button_lamp(BUTTON_COMMAND, head_order.Floor, 0)
		delete_order <- head_order
	}
}

func Elevator_stop(state chan State_t) {
	state <- STOPS
	fmt.Println("The elevator has stopped!\n1. If you wish to order a new floor, do so, or.\n2. Press Ctrl + c to exit program.\n")
	Elevator_clear_all_lights()
	Elev_set_stop_lamp(1)
	time.Sleep(1000 * time.Millisecond)
}

func Elevator_stop_obstruction(head_order_c chan Order, head_order Order, state chan State_t) {
	state <- STOP_OBS
	for {
		time.Sleep(25 * time.Millisecond)
		if !Elev_get_obstruction_signal() {
			head_order_c <- head_order
			return
		}
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
