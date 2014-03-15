package network

import (
	driver "../driver"
	"fmt"
	//"time"
	"net"
	"sort"
)

func priorityHandler(external driver.Client, order_from_cost chan driver.Client, all_clients map[string]driver.Client) {
	fmt.Println("Cost running \n")
	fmt.Println("all clients : ", all_clients)
	for key, value := range all_clients {
		/*if value.Current_floor == external.Floor {
			fmt.Println("er i riktig etg \n")
			value.Ip_from_cost := localIP
			order_from_cost <- external
			flag = true
			break
		}*/
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
	order_from_cost <- designated_client
	fmt.Println("End of cost function \n")
}

func priorityHandler_getCost(client driver.Client, external driver.Client) int {
	cost := 0
	diff := external.Floor - client.Current_floor
	cost = abs(diff)
	direction := client.Direction
	ordered_direction := 0
	if diff < 0 {
		ordered_direction = 1
	} else {
		ordered_direction = -1
	}

	if ordered_direction != direction {
		cost += 5
		for i := 0; i < 4; i++ {
			if client.Order_list[2][i] {
				cost += 1
			}
		}
	}
	return cost
}

func priorityHandler_sort_all_ips(all_clients map[string]driver.Client) net.IP {
	cost_m := make(map[int]net.IP)
	var cost []int
	counter := 0
	for _, value := range all_clients {
		cost_m[value.Cost] = value.Ip
	}
	for i := range cost_m {
		cost = append(cost, i)
		counter++
	}
	sort.Ints(cost)
	fmt.Println("cost map =", cost_m)
	return cost_m[cost[0]]
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
