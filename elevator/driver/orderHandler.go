package driver

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	Ip            net.IP
	Ip_from_cost  net.IP
	Floor         int
	Direction     int
	Button        elev_button_type_t
	Current_floor int
	State         int
	Order_list    [3][4]bool
	Cost          int

}

type Order struct {
	Floor int
	Dir   int
}

func OrderHandler_search_for_orders(order_internal chan Client) {
	var new_order Client
	for {
		time.Sleep(25 * time.Millisecond)
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				if !new_order.Order_list[2][i] {
					new_order.Order_list[2][i] = true
					new_order.Button = BUTTON_COMMAND
					new_order.Floor = i
					Elev_set_button_lamp(BUTTON_COMMAND, i, 1)
					//fmt.Println("Button of the type ",new_order.Button, "pressed. \n")
					fmt.Println("Command order to floor:", i+1, "registered.\n")
					order_internal <- new_order
				}
			}
			if i > 0 {
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					if !new_order.Order_list[1][i] {
						new_order.Floor = i
						new_order.Button = BUTTON_CALL_DOWN
						Elev_set_button_lamp(new_order.Button, new_order.Floor, 1)
						order_internal <- new_order
					}
				}
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					if !new_order.Order_list[0][i] {
						new_order.Floor = i
						new_order.Button = BUTTON_CALL_UP
						Elev_set_button_lamp(new_order.Button, new_order.Floor, 1)
						order_internal <- new_order
					}
				}
			}
		}
	}
}

func OrderHandler_process_orders(order_from_network chan Client, order_to_network chan Client, order_internal chan Client, localIP net.IP) {
	order_internal := make(chan Client)
	head_order_c := make(chan Order)
	prev_order_c := make(chan Order)
	var Head_order Order
	var Prev_order Order
	go Elevator_eventHandler(head_order_c, prev_order_c)
	go OrderHandler_search_for_orders(order_internal)
	for {
		time.Sleep(25 * time.Millisecond)
		select {
		case to_network := <-order_internal:
			fmt.Println("Sending the order on a channel to the network. \n")
			to_network.Ip = localIP
			ifconfi.Current_floor = Prev_order.Floor
			order_to_network <- to_network
		case from_network := <-order_from_network:
			fmt.Println("got order , setting head")
			Head_order = OrderHandler_set_head_order(from_network, Head_order, Prev_order)
			head_order_c <- Head_order
			fmt.Println("from network: ", from_network.Floor+1, "\n")
		case Update_prev := <- prev_order_c:
			Prev_order = Update_prev
		}
	}
}

func Init_orderlist(client Client) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			client.Order_list[i][j] = false
		}
	}
}

func OrderHandler_set_head_order(client Client, Head_order Order, Prev_order Order) Order {
	for {
		time.Sleep(25 * time.Millisecond)
		switch Prev_order.Dir {
		case 1:
			new_head := OrderHandler_case_up(client, Head_order, Prev_order)
			if new_head.Floor != -1 {
				Head_order = new_head
				break
			}
		case -1:
			new_head := OrderHandler_case_down(client, Head_order, Prev_order)
			if new_head.Floor != -1 {
				Head_order = new_head
				break
			}
		}
	}
	return Head_order

}

func OrderHandler_case_up(client Client, Head_order Order, Prev_order Order) Order {
	for i := Prev_order.Floor; i < N_FLOORS; i++ {
		if client.Order_list[0][i] || client.Order_list[1][i] || client.Order_list[2][i] {
			Head_order.Floor = i
			Head_order.Dir = 1
			return Head_order
		}
	}
	Head_order.Floor = -1
	return Head_order
}

func OrderHandler_case_down(client Client, Head_order Order, Prev_order Order) Order {
	for i := Prev_order.Floor; i >= 0; i-- {
		if client.Order_list[0][i] || client.Order_list[1][i] || client.Order_list[2][i] {
			Head_order.Floor = i
			Head_order.Dir = -1
			return Head_order
		}
	}
	Head_order.Floor = -1
	return Head_order
}
