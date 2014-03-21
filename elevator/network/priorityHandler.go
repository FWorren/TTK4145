package network

import (
	driver "../driver"
	"fmt"
	//"time"
	"net"
	"sort"
)

func priorityHandler(external driver.Client, order_from_cost chan driver.Client, all_clients map[string]driver.Client) {
	for key, value := range all_clients {
		value.Cost = priorityHandler_getCost(value, external)
		all_clients[key] = value
	}

	ip := priorityHandler_sort_all_ips(all_clients)
	fmt.Println("order originated from IP :", external.Ip.String())
	designated_client := all_clients[ip.String()]
	designated_client.Floor = external.Floor
	designated_client.Button = external.Button
	designated_client.Ip_from_cost = ip
	fmt.Println("designated client ip :", ip.String())
	fmt.Println("Order at floor: ", designated_client.Floor+1)
	order_from_cost <- designated_client
}

func priorityHandler_getCost(client driver.Client, external driver.Client) int {
	cost := 0
	if client.State == driver.STOPS || client.State == driver.STOP_OBS {
		cost += 20
	}
	diff := external.Floor - client.Current_floor
	cost = abs(diff)

	for i := 0; i < driver.N_FLOORS; i++ {
		if client.Order_list[driver.BUTTON_COMMAND][i] {
			cost += 2
		}
		if client.Order_list[driver.BUTTON_CALL_DOWN][i] {
			cost += 2
		}
		if client.Order_list[driver.BUTTON_CALL_UP][i] {
			cost += 2
		}
	}

	/*if diff <0 && external.Button == driver.BUTTON_CALL_UP {
			for i := client.Current_floor-1; i >= 0; i-- {
				if client.Order_list[driver.BUTTON_CALL_DOWN][i] {
					cost += 1
				}
			}
		}else if diff >0 && external.Button == driver.BUTTON_CALL_DOWN {
			for i := 0; i < ; i++ {

			}
		}

	/*if diff > 0 {
		for i := client.Current_floor + 1; i >= 0; i-- {
			if client.Order_list[driver.BUTTON_CALL_UP][i] {
				cost += 1
			}
		}
	} else if diff < 0 {
		for i := client.Current_floor - 1; i < 4; i++ {
			if client.Order_list[driver.BUTTON_CALL_DOWN][i] {
				cost += 1
			}
		}*/

	return cost
}

func priorityHandler_sort_all_ips(all_clients map[string]driver.Client) net.IP {
	cost_m := make(map[int]net.IP)
	var cost []int
	for _, value := range all_clients {
		cost_m[value.Cost] = value.Ip
	}
	for key, _ := range cost_m {
		cost = append(cost, key)
	}
	sort.Ints(cost)
	fmt.Println("cost sorted: ", cost)
	return cost_m[cost[0]]
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
