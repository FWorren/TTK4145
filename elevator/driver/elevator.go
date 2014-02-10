package driver

import (
	"fmt"
)

type State_t int

const  (
    RUN State_t = iota
    WAIT
    DOOR
    STOPS
    UNDEF
    EXIT
)

var State State_t

func Elevator_init() {
	OrderLogic_init()
	Elevator_clear_all_lights()
	if Elev_get_floor_sensor_signal() != -1 {
		State = WAIT
		OrderLogic_update_previous_order(Elev_get_floor_sensor_signal(),-1)
		Elev_set_floor_indicator(Previous_order.floor)
	}
	Elev_set_speed(-300)
	for {
		if Elev_get_floor_sensor_signal() != -1 {
			State = WAIT
			OrderLogic_update_previous_order(Elev_get_floor_sensor_signal(),-1)
			Elev_set_floor_indicator(Previous_order.floor)
			Elev_set_speed(0)
			break
		}
	}
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
			case UNDEF:
				Elevator_init()
			default:
				fmt.Println("Error! Program terminated!\n")
				State = EXIT
		}
		if State == EXIT {break}
	}
}

func Elevator_wait() {
	for {
		OrderLogic_search_for_orders()
		if OrderLogic_get_number_of_internal_orders() > 0 {State = RUN}
		if Elev_get_stop_signal() {State = STOPS}
		if State != WAIT{break}
	}
}

func Elevator_run() {
	OrderLogic_set_head_order();
	Elev_set_speed(300*Head_order.dir);
	for {
		OrderLogic_search_for_orders()
		floor := Elev_get_floor_sensor_signal()
		if floor != -1 {
			OrderLogic_update_previous_order(floor,Head_order.dir)
			Elev_set_floor_indicator(floor)
		}
		if floor == Head_order.floor {
			State = DOOR
		}
		if Elev_get_stop_signal() || Elev_get_obstruction_signal() {State = STOPS}
		if State != RUN{break}
	}
}

func Elevator_door() {
	for {
		if State != DOOR{break}
	}
}

func Elevator_stop() {
	for {
		if State != STOPS{break}
	}
}

func Elevator_clear_lights_current_floor(current_floor int) {
	Elev_set_button_lamp(BUTTON_COMMAND,current_floor,0)
	if current_floor > 0 {
		Elev_set_button_lamp(BUTTON_CALL_DOWN,current_floor,0)
	}
	if current_floor < N_FLOORS - 1 {
		Elev_set_button_lamp(BUTTON_CALL_UP,current_floor,0)
	}
}

func Elevator_clear_all_lights() {
	Elev_set_door_open_lamp(0)
	Elev_set_stop_lamp(0)
	for i := 0; i < N_FLOORS; i++ {
		Elev_set_button_lamp(BUTTON_COMMAND,i,0)
		if i > 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN,i,0)
		}
		if i < N_FLOORS - 1 {
			Elev_set_button_lamp(BUTTON_CALL_UP,i,0)
		}
	}
}

/*func Elevator_break(direction int) {
	Elev_set_speed(20*(-direction))
	Elevator_timer(2)
	Elev_set_speed(0)
}

func Elevator_timer(sec int) {
	time.Sleep(sec*int(time.Millisecond))
}*/