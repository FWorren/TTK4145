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
	fmt.Println("Order at floor: ", designated_client.Floor+1, "\n")
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
	return cost
}

func priorityHandler_sort_all_ips(all_clients map[string]driver.Client) net.IP {
	cost_m := make(map[string]int)
	var designated_ip net.IP
	var cost []int
	counter := 0
	for _, value := range all_clients {
		cost_m[value.Ip.String()] = value.Cost
	}
	for _, value := range cost_m {
		cost = append(cost, value)
		counter++
	}
	sort.Ints(cost)

	fmt.Println("cost sorted: ", cost)
	for key, value := range cost_m {
		if value == cost[0] {
			for key_ip, value_c := range all_clients {
				if key_ip == key {
					designated_ip = value_c.Ip
				}
			}
		}
	}
	return designated_ip
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
