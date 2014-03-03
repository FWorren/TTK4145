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

var Order_list = [3][4]bool{
	    {false, false, false, false},
	    {false, false, false, false},
	    {false, false, false, false},
}

func OrderHandler_search_for_orders(order_internal chan Client) {
	var new_order Client
	for {
		time.Sleep(25*time.Millisecond)
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				if !Order_list[2][i]{
					Order_list[2][i]= true
					Elev_set_button_lamp(BUTTON_COMMAND,i,1)
					fmt.Println("Command order: ",i+1)
					
				}
			}
			if i > 0{
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					if !Order_list[1][i]{
						Order_list[1][i]= true
						new_order.Floor = i
						new_order.Button = BUTTON_CALL_DOWN
						Elev_set_button_lamp(new_order.Button,new_order.Floor,1)
						order_internal <- new_order
					}
				}	
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					if !Order_list[0][i]{
						Order_list[0][i]= true
						new_order.Floor = i
						new_order.Button = BUTTON_CALL_UP
						Elev_set_button_lamp(new_order.Button,new_order.Floor,1)
						order_internal <- new_order
					}
				}	
			}
			
		}

	}
}

func OrderHandler_process_orders(order_to_network chan Client, order_from_network chan Client,order_internal chan Client) {
	go OrderHandler_search_for_orders(order_internal)
	for {
		time.Sleep(25*time.Millisecond)
		select {
			case to_network := <- order_internal:
				order_to_network <- to_network
			case from_network := <- order_from_network:
				Order_list[from_network.Button][from_network.Floor]= true
				fmt.Println("from network: ", from_network.Floor +1,"\n")	
		}
	}
}