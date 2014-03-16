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
	Button        Elev_button_type_t
	Current_floor int
	State         State_t
	Order_list    [3][4]bool
	Cost          int
}

type Order struct {
	Floor  int
	Dir    int
	Button Elev_button_type_t
}

type Lights struct {
	Floor  int
	Button Elev_button_type_t
}

func OrderHandler_process_orders(order_from_network chan Client, order_to_network chan Client, current_floor Order, localIP net.IP) {
	order_internal := make(chan Order)
	head_order_c := make(chan Order)
	prev_order_c := make(chan Order)
	del_Order := make(chan Order)
	//convenient_floor := make(chan bool)
	numb_orders_c := make(chan int)
	state_c := make(chan State_t)
	local_list_c := make(chan [3][4]bool)

	var state State_t
	var local_list [3][4]bool
	var client Client
	var Head_order Order

	state = UNDEF
	Prev_order := current_floor

	go Elevator_eventHandler(head_order_c, prev_order_c, del_Order, state_c)
	go OrderHandler_search_for_orders(order_internal, local_list)
	go Check_number_of_local_orders(numb_orders_c, local_list_c)

	for {
		time.Sleep(25 * time.Millisecond)
		select {
		case to_network := <-order_internal:
			if to_network.Button == BUTTON_COMMAND {
				if !local_list[BUTTON_COMMAND][to_network.Floor] {
					local_list[BUTTON_COMMAND][to_network.Floor] = true
					client.Order_list[BUTTON_COMMAND][to_network.Floor] = true
					local_list_c <- local_list
				}
			} else {
				if client.State != STOP_OBS && !local_list[to_network.Button][to_network.Floor]{
					fmt.Println("Sending the order on a channel to the network. \n")
					client.Floor = to_network.Floor
					client.Button = to_network.Button
					client.Ip = localIP
					client.Current_floor = Prev_order.Floor
					order_to_network <- client
				}
			}
		case from_network := <-order_from_network:
			fmt.Println("recieving from network")
			if !local_list[from_network.Button][from_network.Floor] {
				local_list[from_network.Button][from_network.Floor] = true
				client.Order_list[from_network.Button][from_network.Floor] = true
				local_list_c <- local_list
			}
		case state = <-state_c:
			fmt.Println("state = ", state)
			client.State = state
		case numb_orders := <-numb_orders_c:
			fmt.Println("Number of orders: ", numb_orders)
			if (state == WAIT || state == UNDEF) && numb_orders > 0 {
				Head_order = OrderHandler_set_head_order(local_list, Head_order, Prev_order)
				head_order_c <- Head_order
			}
		case Update_prev := <-prev_order_c:
			Prev_order = Update_prev
			client.Direction = Prev_order.Dir
		case del_msg := <-del_Order:
			local_list[del_msg.Button][del_msg.Floor] = false
			client.Order_list[del_msg.Button][del_msg.Floor] = false
			local_list_c <- local_list
		}
	}
}

func OrderHandler_search_for_orders(order_internal chan Order, local_list [3][4]bool) {
	var new_order Order
	for {
		time.Sleep(25 * time.Millisecond)
		for i := 0; i < N_FLOORS; i++ {
			if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
				new_order.Button = BUTTON_COMMAND
				new_order.Floor = i
				Elev_set_button_lamp(BUTTON_COMMAND, i, 1)
				order_internal <- new_order
			}
			if i > 0 {
				if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
					if !local_list[1][i] {
						new_order.Floor = i
						new_order.Button = BUTTON_CALL_DOWN
						order_internal <- new_order
					}
				}
			}
			if i < N_FLOORS-1 {
				if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
					if !local_list[0][i] {
						new_order.Floor = i
						new_order.Button = BUTTON_CALL_UP
						order_internal <- new_order
					}
				}
			}
		}
	}
}

/*func Init_orderlist(client Client) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			client.Order_list[i][j] = false
		}
	}
}*/

func Check_number_of_local_orders(numb_orders_c chan int, local_list_c chan [3][4]bool) {
	numb_orders := 0
	for {
		time.Sleep(25 * time.Millisecond)
		select {
		case local_list := <-local_list_c:
			numb_orders = 0
			for i := 0; i < N_FLOORS; i++ {
				if local_list[BUTTON_CALL_UP][i] {
					numb_orders++
				}
				if local_list[BUTTON_CALL_DOWN][i] {
					numb_orders++
				}
				if local_list[BUTTON_COMMAND][i] {
					numb_orders++
				}
			}
			if numb_orders > 0 {
				numb_orders_c <- numb_orders
			}
		}
	}
}

func OrderHandler_set_head_order(local_list [3][4]bool, Head_order Order, Prev_order Order) Order {
	for {
		time.Sleep(25 * time.Millisecond)
		switch Prev_order.Dir {
		case 1:
			new_head := OrderHandler_state_up(local_list, Head_order, Prev_order)
			if new_head.Floor != -1 {
				Head_order = new_head
				return Head_order
			}
			Prev_order.Dir = new_head.Dir
		case -1:
			new_head := OrderHandler_state_down(local_list, Head_order, Prev_order)
			if new_head.Floor != -1 {
				Head_order = new_head
				return Head_order
			}
			Prev_order.Dir = new_head.Dir
		}

	}
}

func OrderHandler_state_up(local_list [3][4]bool, Head_order Order, Prev_order Order) Order {
	if Prev_order.Floor == N_FLOORS-1 {
		Head_order.Dir = -1
		Head_order.Floor = -1
		return Head_order
	}
	for i := Prev_order.Floor; i < N_FLOORS; i++ {
		if local_list[0][i] {
			Head_order.Floor = i
			Head_order.Dir = 1
			Head_order.Button = BUTTON_CALL_UP
			return Head_order
		}
		if local_list[1][i] {
			Head_order.Floor = i
			Head_order.Dir = 1
			Head_order.Button = BUTTON_CALL_DOWN
			return Head_order
		}
		if local_list[2][i] {
			Head_order.Floor = i
			Head_order.Dir = 1
			Head_order.Button = BUTTON_COMMAND
			return Head_order
		}
	}
	Head_order.Floor = -1
	Head_order.Dir = -1
	return Head_order
}

func OrderHandler_state_down(local_list [3][4]bool, Head_order Order, Prev_order Order) Order {
	if Prev_order.Floor == 0 {
		Head_order.Dir = 1
		Head_order.Floor = -1
		return Head_order
	}
	for i := Prev_order.Floor; i >= 0; i-- {
		if local_list[0][i] {
			Head_order.Floor = i
			Head_order.Dir = -1
			Head_order.Button = BUTTON_CALL_UP
			return Head_order
		}
		if local_list[1][i] {
			Head_order.Floor = i
			Head_order.Dir = -1
			Head_order.Button = BUTTON_CALL_DOWN
			return Head_order
		}
		if local_list[2][i] {
			Head_order.Floor = i
			Head_order.Dir = -1
			Head_order.Button = BUTTON_COMMAND
			return Head_order
		}
	}
	Head_order.Floor = -1
	Head_order.Dir = 1
	return Head_order
}

/*
func OrderHandler_check_convenient_order(local_list [3][4]bool, Prev_order Order) bool {
	dir := Prev_order.Dir
	floor := Prev_order.Floor
	if dir == -1 {
		if local_list[BUTTON_CALL_DOWN][floor] || local_list[BUTTON_COMMAND][floor] {
			return true
		}
	} else {
		if local_list[BUTTON_CALL_UP][floor] || local_list[BUTTON_COMMAND][floor] {
			return true
		}
	}
	return false
}*/
