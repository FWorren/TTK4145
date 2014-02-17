package driver

import (
	"fmt"
	"encoding/json"
	//"time"
)

var command [N_FLOORS]int

type Order_state_t int

const (
	UP Order_state_t = iota
	DOWN
	SET
)

var Order_state Order_state_t

type Order struct {
	floor int
	dir   int
}

var Head_order Order
var Previous_order Order

func OrderLogic_init() {
	for i := 0; i < N_FLOORS; i++ {
		command[i] = 0
	}
}

/*
func OrderLogic_search_for_orders() {
	for {
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				if command[i] != 1 {
					command[i] = 1
					Elev_set_button_lamp(BUTTON_COMMAND, i, 1)
				}
			}
		}
	}
}*/

func OrderLogic_search_for_orders(order_internal chan []byte) {
	type Order_t struct {
		floor int
		button elev_button_type_t
	}
	var New_order Order_t 
	for {
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				New_order.floor = i
				New_order.button = BUTTON_COMMAND
				b, err := json.Marshal(New_order)
				order_internal <- b
			}
			if i > 0 {
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					New_order.floor = i
					New_order.button = BUTTON_CALL_DOWN
					b, err := json.Marshal(New_order)
					order_internal <- b
				}	
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					New_order.floor = i
					New_order.button = BUTTON_CALL_UP
					b, err := json.Marshal(New_order)
					order_internal <- b
				}	
			}
		}
	}
}

func OrderLogic_set_order(order_internal chan []byte, order_from_network chan int) {
	for {
		select {
			case internal := <- order_internal:
				var msg struct{}
				err := json.Unmarshal(internal,&msg)
				Elev_set_button_lamp(msg.button,msg.floor,1)
			case network := <- order_from_network:
				command[0] = 1
			default:

		}
	}
}

func OrderLogic_set_head_order() {
	OrderLogic_set_order_state()
	for {
		switch Order_state {
			case UP:
				OrderLogic_state_up()
			case DOWN:
				OrderLogic_state_down()
			default:
				fmt.Println("Error! No queue!\n")
				Order_state = SET
		}
		if Order_state == SET {
			break
		}
	}
}

func OrderLogic_set_order_state() {
	if Previous_order.dir == 1 {
		Order_state = UP
	} else {
		Order_state = DOWN
	}
}

func OrderLogic_state_up() {
	if Previous_order.floor == N_FLOORS-1 {
		Order_state = DOWN
		return
	}
	for i := Previous_order.floor + 1; i < N_FLOORS; i++ {
		if command[i] == 1 {
			Head_order.floor = i
			Head_order.dir = 1
			Order_state = SET
			return
		}
	}
	Order_state = DOWN
}

func OrderLogic_state_down() {
	if Previous_order.floor == 0 {
		Order_state = UP
		return
	}
	for i := Previous_order.floor - 1; i >= 0; i-- {
		if command[i] == 1 {
			Head_order.floor = i
			Head_order.dir = -1
			Order_state = SET
			return
		}
	}
	Order_state = UP
}

func OrderLogic_get_number_of_orders() int {
	n_orders := 0
	for i := 0; i < N_FLOORS; i++ {
		if command[i] == 1 {
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
}
