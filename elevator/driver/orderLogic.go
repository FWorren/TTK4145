package driver

import (
	//"fmt"
	//"encoding/json"
	//"time"
)

var command [N_FLOORS]int
var Up[N_FLOORS]int
var Down[N_FLOORS]int

type State_dir_t int

const (
	UP State_dir_t = iota
	DOWN
	SET
)

type Order_t struct {
	floor int
	dir   int
}

type Order_c struct {
	floor int
	button elev_button_type_t
}

var Head_order Order_t
var Previous_order Order_t

func OrderLogic_init() {
	for i := 0; i < N_FLOORS; i++ {
		command[i] = 0
		Up[i] = 0
		Down[i] = 0
	}
}

func OrderLogic_search_for_orders(order_internal chan Order_c) {
	var new_order Order_c
	for {
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				new_order.floor = i
				new_order.button = BUTTON_COMMAND
				order_internal <- new_order
			}
			if i > 0{
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					new_order.floor = i
					new_order.button = BUTTON_CALL_DOWN
					order_internal <- new_order
				}	
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					new_order.floor = i
					new_order.button = BUTTON_CALL_UP
					order_internal <- new_order
				}	
			}
		}
	}
}

func OrderLogic_set_order(order_internal chan Order_c, order_from_network chan Order_c) {
	for {
		select {
			case internal := <- order_internal:
				Elev_set_button_lamp(internal.button,internal.floor,1)
				switch internal.button {
					case BUTTON_COMMAND:
						command[internal.floor] = 1
					case BUTTON_CALL_UP:
						Up[internal.floor] = 1
					case BUTTON_CALL_DOWN:
						Down[internal.floor] = 1
				}
			/*case from_network := <- order_from_network:
				var msg []Order_c
				err := json.Unmarshal([]byte(from_network),&msg)
				if err != nil {
					fmt.Println("error: ", err)	
				}
				fmt.Println(msg)*/
		}
	}
}

func OrderLogic_set_head_order() {
	var State_dir State_dir_t
	State_dir = OrderLogic_set_order_state()
	for {
		switch State_dir {
			case UP:
				State_dir = OrderLogic_state_up()
			case DOWN:
				State_dir = OrderLogic_state_down()
			case SET:
				return
		}
	}
}

func OrderLogic_set_order_state() State_dir_t {
	if Previous_order.dir == 1 {
		return UP
	} else {
		return DOWN
	}
}

func OrderLogic_state_up() State_dir_t {
	if Previous_order.floor == N_FLOORS-1 {
		return DOWN
	}
	for i := Previous_order.floor + 1; i < N_FLOORS; i++ {
		if command[i] == 1 || Up[i] == 1 || Down[i] == 1 {
			Head_order.floor = i
			Head_order.dir = 1
			return SET
		}
	}
	return DOWN
}

func OrderLogic_state_down() State_dir_t {
	if Previous_order.floor == 0 {
		return UP
	}
	for i := Previous_order.floor - 1; i >= 0; i-- {
		if command[i] == 1 || Up[i] == 1 || Down[i] == 1 {
			Head_order.floor = i
			Head_order.dir = -1
			return SET
		}
	}
	return UP
}

func OrderLogic_get_number_of_orders() int {
	n_orders := 0
	for i := 0; i < N_FLOORS; i++ {
		if command[i] == 1 || Up[i] == 1 || Down[i] == 1 {
			n_orders += 1
		}
	}
	return n_orders
}

func OrderLogic_delete_all_orders() {
	for i := 0; i < N_FLOORS; i++ {
		command[i] = 0
	}
	Head_order.floor = 0
	Head_order.dir = 0
}

func OrderLogic_update_previous_order(floor int, direction int) {
	Previous_order.floor = floor
	Previous_order.dir = direction
}

func OrderLogic_delete_current_order(current_floor int) {
	Previous_order.floor = Head_order.floor
	Previous_order.dir = Head_order.dir
	if command[current_floor] == 1 {
		command[current_floor] = 0
	}
	if current_floor > 0 && Down[current_floor] == 1 {
		Down[current_floor] = 0
	}
	if current_floor < N_FLOORS-1 && Up[current_floor] == 1 {
		Up[current_floor] = 0
	}
}