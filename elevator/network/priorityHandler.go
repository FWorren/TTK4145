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
	flag := false
	for _, value := range all_clients {
		if value.Current_floor == external.Floor {
			fmt.Println("er i riktig etg \n")
			order_from_cost <- external
			flag = true
			break
		}
		value.Cost = priorityHandler_getCost(value, external)
	}
	if flag {
		return
	} else {
		ip := priorityHandler_sort_all_ips(all_clients)
		designated_client := all_clients[ip.String()]
		fmt.Println("Order originated from IP :",designated_client.Ip.String(),"\n")
		designated_client.Ip_from_cost = ip
		fmt.Println("The designated client IP :",designated_client.Ip_from_cost.String(),"\n")
		designated_client.Order_list[external.Button][external.Floor] = true
		fmt.Println("Orderlist : ", designated_client.Order_list, "\n")
		order_from_cost <- designated_client
	}
	fmt.Println("End of cost function \n")
}

func priorityHandler_getCost(client driver.Client, external driver.Client) int {
	cost := abs(client.Current_floor - external.Floor)
	for i := 0; i < 4; i++ {
		if client.Order_list[2][i] {
			cost += 1
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
	//fmt.Println("Counter = ", counter)
	sort.Ints(cost)

	//fmt.Println(cost_m[cost[0]], "\n")
	return cost_m[cost[0]]
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
