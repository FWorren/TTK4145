package driver

import (
	//"fmt"
	"net"
	"time"
)

type NetState_t int

const (
	OFF NetState_t = iota
	ON
)

type Client struct {
	Ip            net.IP
	Ip_from_cost  net.IP
	Floor         int
	Direction     int
	Button        Elev_button_type_t
	Current_floor int
	State         State_t
	NetState      NetState_t
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
	Flag   bool
}

func OrderHandler_process_orders(order_from_network chan Client, order_to_network chan Client, check_backup_c chan Client, status_update_c chan Client, send_lights_c chan Lights, send_del_req_c chan Order, order_complete_c chan Order, disconnected chan int, netstate_c chan NetState_t, current_floor Order, localIP net.IP) {
	order_internal := make(chan Order, 1)
	head_order_c := make(chan Order, 1)
	prev_order_c := make(chan Order, 1)
	del_Order := make(chan Order, 1)
	reset_list_c := make(chan Order, 1)
	state_c := make(chan State_t, 1)
	reset_all_c := make(chan int, 1)
	convenient_list_c := make(chan [3][4]bool, 1)

	var state State_t
	var netState NetState_t
	var local_list [3][4]bool
	var client Client
	var Head_order Order
	var light Lights
	state = UNDEF
	Prev_order := current_floor
	client.Current_floor = current_floor.Floor
	client.Direction = current_floor.Dir

	go Elevator_eventHandler(head_order_c, prev_order_c, del_Order, state_c, convenient_list_c)
	go OrderHandler_search_for_orders(order_internal, reset_list_c, reset_all_c)
	timeOut := make(<-chan time.Time)
	for {
		timeOut = time.After(100 * time.Millisecond)
		select {
		case to_network := <-order_internal:
			client.Floor = to_network.Floor
			client.Button = to_network.Button
			client.Ip = localIP
			client.Current_floor = Prev_order.Floor
			if to_network.Button == BUTTON_COMMAND {
				if !local_list[BUTTON_COMMAND][to_network.Floor] {
					local_list[BUTTON_COMMAND][to_network.Floor] = true
					client.Order_list[BUTTON_COMMAND][to_network.Floor] = true
					if netState == ON {
						order_to_network <- client
					}
				}
			} else {
				if netState == ON {
					order_to_network <- client
				} else {
					reset_list_c <- to_network
				}
			}

		case from_network := <-order_from_network:
			local_list[from_network.Button][from_network.Floor] = true
			client.Order_list[from_network.Button][from_network.Floor] = true
			light.Floor = from_network.Floor
			light.Button = from_network.Button
			light.Flag = true
			send_lights_c <- light
		case backup_client := <-check_backup_c:
			client.Order_list = get_backup_orders(backup_client)
			local_list = client.Order_list
		case state = <-state_c:
			client.State = state
		case Update_prev := <-prev_order_c:
			Prev_order = Update_prev
			client.Direction = Prev_order.Dir
			client.Current_floor = Prev_order.Floor
			convenient_list_c <- local_list
		case del_msg := <-del_Order:
			local_list[del_msg.Button][del_msg.Floor] = false
			client.Order_list[del_msg.Button][del_msg.Floor] = false
			reset_list_c <- del_msg
			light.Floor = del_msg.Floor
			light.Button = del_msg.Button
			light.Flag = false
			send_lights_c <- light
			send_del_req_c <- del_msg
		case completed_order := <-order_complete_c:
			reset_list_c <- completed_order
		case netState = <-netstate_c:
			client.NetState = netState
		case <-disconnected:
			for i := 0; i < N_FLOORS; i++ {
				local_list[BUTTON_CALL_DOWN][i] = false
				local_list[BUTTON_CALL_UP][i] = false
				client.Order_list[BUTTON_CALL_DOWN][i] = false
				client.Order_list[BUTTON_CALL_UP][i] = false
				Elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0)
				Elev_set_button_lamp(BUTTON_CALL_UP, i, 0)
			}
			reset_all_c <- 1
		case <-timeOut:
			has_order := Check_number_of_local_orders(local_list)
			if (state == WAIT || state == UNDEF) && has_order {
				Head_order = OrderHandler_set_head_order(local_list, Head_order, Prev_order)
				client.Direction = Head_order.Dir
				head_order_c <- Head_order
			}
			status_update_c <- client
		}
	}
}

func get_backup_orders(client Client) [3][4]bool {
	var command_list [3][4]bool
	for i := 0; i < N_FLOORS; i++ {
		if client.Order_list[BUTTON_COMMAND][i] {
			command_list[BUTTON_COMMAND][i] = true
			Elev_set_button_lamp(BUTTON_COMMAND, i, 1)
		}
	}
	return command_list
}

func OrderHandler_search_for_orders(order_internal chan Order, reset_list_c chan Order, reset_all_c chan int) {
	var new_order Order
	var list [3][4]bool

	go func() {
		for {
			select {
			case reset_floor := <-reset_list_c:
				list[reset_floor.Button][reset_floor.Floor] = false
			case <-reset_all_c:
				for i := 0; i < N_FLOORS; i++ {
					list[BUTTON_CALL_DOWN][i] = false
					list[BUTTON_CALL_UP][i] = false
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			for i := 0; i < N_FLOORS; i++ {
				if Elev_get_button_signal(BUTTON_COMMAND, i) == 1 {
					if !list[BUTTON_COMMAND][i] {
						new_order.Button = BUTTON_COMMAND
						new_order.Floor = i
						list[BUTTON_COMMAND][i] = true
						Elev_set_button_lamp(BUTTON_COMMAND, i, 1)
						order_internal <- new_order
					}
				}
				if i > 0 {
					if Elev_get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
						if !list[BUTTON_CALL_DOWN][i] {
							new_order.Floor = i
							new_order.Button = BUTTON_CALL_DOWN
							list[BUTTON_CALL_DOWN][i] = true
							order_internal <- new_order
						}
					}
				}
				if i < N_FLOORS-1 {
					if Elev_get_button_signal(BUTTON_CALL_UP, i) == 1 {
						if !list[BUTTON_CALL_UP][i] {
							new_order.Floor = i
							new_order.Button = BUTTON_CALL_UP
							list[BUTTON_CALL_UP][i] = true
							order_internal <- new_order
						}
					}
				}
			}
		}
	}()
}

func Check_number_of_local_orders(local_list [3][4]bool) bool {
	numb_orders := 0
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
		return true
	} else {
		return false
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
		if local_list[BUTTON_CALL_UP][i] {
			Head_order.Floor = i
			Head_order.Dir = 1
			Head_order.Button = BUTTON_CALL_UP
			return Head_order
		}
		if local_list[BUTTON_CALL_DOWN][i] {
			Head_order.Floor = i
			Head_order.Dir = 1
			Head_order.Button = BUTTON_CALL_DOWN
			return Head_order
		}
		if local_list[BUTTON_COMMAND][i] {
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
		if local_list[BUTTON_CALL_UP][i] {
			Head_order.Floor = i
			Head_order.Dir = -1
			Head_order.Button = BUTTON_CALL_UP
			return Head_order
		}
		if local_list[BUTTON_CALL_DOWN][i] {
			Head_order.Floor = i
			Head_order.Dir = -1
			Head_order.Button = BUTTON_CALL_DOWN
			return Head_order
		}
		if local_list[BUTTON_COMMAND][i] {
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

func OrderHandler_check_convenient_order(local_list [3][4]bool, Prev_order Order) (bool, Order) {
	var Convenient_order Order
	flag := false
	dir := Prev_order.Dir
	floor := Prev_order.Floor
	if dir == -1 {
		if local_list[BUTTON_CALL_DOWN][floor] {
			Convenient_order.Button = BUTTON_CALL_DOWN
			Convenient_order.Floor = floor
			Convenient_order.Dir = dir
			flag = true
			return flag, Convenient_order
		} else if local_list[BUTTON_COMMAND][floor] {
			Convenient_order.Button = BUTTON_COMMAND
			Convenient_order.Floor = floor
			Convenient_order.Dir = dir
			flag = true
			return flag, Convenient_order
		}
	}
	if dir == 1 {
		if local_list[BUTTON_CALL_UP][floor] {
			Convenient_order.Button = BUTTON_CALL_UP
			Convenient_order.Floor = floor
			Convenient_order.Dir = dir
			flag = true
			return flag, Convenient_order
		} else if local_list[BUTTON_COMMAND][floor] {
			Convenient_order.Button = BUTTON_COMMAND
			Convenient_order.Floor = floor
			Convenient_order.Dir = dir
			flag = true
			return flag, Convenient_order
		}
	}
	return flag, Convenient_order
}
