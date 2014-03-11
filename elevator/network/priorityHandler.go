package network

import (
	driver "../driver"
	"fmt"
	//"time"
	"sort"
)

func priorityHandler(external driver.Client, order_from_cost chan driver.Client, all_clients map[string]driver.Client) {
	fmt.Println("Cost running \n")
	flag := false
	for _, value := range all_clients {
		if value.Current_floor == external.Floor {
			order_from_cost <- external
			flag = true
			break
		}
		value.Cost = priorityHandler_getCost(value, external)
	}
	if flag {
		return
	} else {
		sorted_ips := priorityHandler_sort_all_ips(all_clients)
		designated_client := all_clients[sorted_ips[0]]
		designated_client.Order_list[external.Button][external.Floor] = true
		order_from_cost <- all_clients[sorted_ips[0]]
	}

	/*switch external.Button {
		case driver.BUTTON_CALL_UP:
			fmt.Println("Order request UPWARD from floor:", external.Floor+1, "\n")
			driver.Elev_set_button_lamp(external.Button, external.Floor, 1)
		case driver.BUTTON_CALL_DOWN:
			fmt.Println("Order request DOWNWARD from floor:", external.Floor+1, "\n")
			driver.Elev_set_button_lamp(external.Button, external.Floor, 1)
		case driver.BUTTON_COMMAND:
			fmt.Println("Order from inside the elevator to floor:", external.Floor+1, "\n")
	}*/
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

func priorityHandler_sort_all_ips(all_clients map[string]driver.Client) []string {
	var cost_m map[int]string
	var cost []int
	counter := 0
	for i := range cost_m {
		cost = append(cost, i)
		counter++
	}
	sort.Ints(cost)
	var ip_m = make([]string, counter)
	for i := 0; i < counter; i++ {
		ip_m[i] = cost_m[cost[i]]
	}

	return ip_m
}

func abs(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}
