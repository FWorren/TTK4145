package driver

import (
	"fmt"
	"net"
	//"encoding/json"
	"time"
)

type Client struct {
	Ip net.IP
	Floor int
	Button elev_button_type_t
}

func OrderLogic_search_for_orders(order_internal chan Client) {
	var new_order Client
	for {
		time.Sleep(25*time.Millisecond)
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				new_order.Floor = i
				new_order.Button = BUTTON_COMMAND
				Elev_set_button_lamp(new_order.Button,new_order.Floor,1)
				order_internal <- new_order
				
			}
			if i > 0{
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					//time.Sleep(25*time.Millisecond)
					new_order.Floor = i
					new_order.Button = BUTTON_CALL_DOWN
					Elev_set_button_lamp(new_order.Button,new_order.Floor,1)
					order_internal <- new_order
					
				}	
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					//time.Sleep(25*time.Millisecond)
					new_order.Floor = i
					new_order.Button = BUTTON_CALL_UP
					Elev_set_button_lamp(new_order.Button,new_order.Floor,1)
					order_internal <- new_order
					
				}	
			}
			
		}

	}
}

func OrderLogic_process_orders(order_to_network chan Client, order_from_network chan Client,order_internal chan Client) {
	for {
		time.Sleep(25*time.Millisecond)
		select {
			case to_network := <- order_internal:
				order_to_network <- to_network
			case from_network := <- order_from_network:
				fmt.Println("from network: ", from_network.Floor,"\n")	
		}
	}
}