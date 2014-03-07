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

func Elevator_init(State State_t) State_t {
	Elevator_clear_all_lights()
	if Elev_get_floor_sensor_signal() != -1 {
		return WAIT
	}
	Elev_set_speed(-300)
	for {
		if Elev_get_floor_sensor_signal() != -1 {
			Elevator_break(-1)
			return WAIT
		}
	}
}

func Elevator_statemachine() {
	var State State_t
	State = UNDEF
	for {
		switch State {
			case RUN:
				State = Elevator_run()
			case WAIT:
				State = Elevator_wait()
			case DOOR:
				State = Elevator_door()
			case STOPS:
				State = Elevator_stop()
			case STOP_OBS:
				State = Elevator_stop_obstruction()
			case UNDEF:
				State = Elevator_init(State)
			case EXIT:
				return
		}
	}
}

func Elevator_wait() State_t {
	for {
		time.Sleep(25*time.Millisecond)
		if Elev_get_stop_signal() {
			return STOPS
		}
	}
}

func Elevator_run() State_t {
	for {
		time.Sleep(25*time.Millisecond)
		//floor := Elev_get_floor_sensor_signal()
	}
}

func Elevator_door() State_t {
	if Elev_get_floor_sensor_signal() != -1 {
		Elev_set_door_open_lamp(1)
	}
	time.Sleep(3 * time.Second)
	for {
		time.Sleep(25*time.Millisecond)
		if !Elev_get_obstruction_signal() {
			return WAIT
		}
	}
}

func Elevator_stop() State_t {
	fmt.Println("The elevator has stopped!\n1. If you wish to order a new floor, do so, or.\n2. Press Ctrl + c to exit program.\n")
	Elevator_clear_all_lights()
	Elev_set_stop_lamp(1)
	for {
		
	}
}

func Elevator_stop_obstruction() State_t {
	for {
		time.Sleep(25*time.Millisecond)
		if !Elev_get_obstruction_signal() {
			return RUN
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
	Elev_set_speed(0)
}