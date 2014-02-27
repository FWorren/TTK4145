package driver

import (
	"fmt"
	"net"
	//"encoding/json"
	//"time"
)

type Client struct {
	ip net.IP
	floor int
	button elev_button_type_t
}

func OrderLogic_search_for_orders(order_to_network chan Client) {
	var new_order Client
	for {
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				new_order.floor = i
				new_order.button = BUTTON_COMMAND
				order_to_network <- new_order
			}
			if i > 0{
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					new_order.floor = i
					new_order.button = BUTTON_CALL_DOWN
					order_to_network <- new_order
				}	
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					new_order.floor = i
					new_order.button = BUTTON_CALL_UP
					order_to_network <- new_order
				}	
			}
		}
	}
}

func OrderLogic_set_order(order_to_network chan Client, order_from_network chan Client) {
	for {
		select {
			case to_network := <- order_to_network:
				switch to_network.button {
					case BUTTON_COMMAND:
						Elev_set_button_lamp(to_network.button,to_network.floor,1)
					case BUTTON_CALL_UP:
						Elev_set_button_lamp(to_network.button,to_network.floor,1)
					case BUTTON_CALL_DOWN:
						Elev_set_button_lamp(to_network.button,to_network.floor,1)
				}
			case from_network := <- order_from_network:
				fmt.Println("error: ", from_network)	
		}
	}
}