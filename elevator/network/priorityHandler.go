package network

import (
	driver "../driver"
	"fmt"
	//"time"
)

func priorityHandler(msg driver.Client, send_from_network chan driver.Client) {
	fmt.Println("Cost running \n")
	switch msg.Button {
		case driver.BUTTON_CALL_UP:
			fmt.Println("Order request UPWARD from floor:", msg.Floor+1, "\n")
			driver.Elev_set_button_lamp(msg.Button, msg.Floor, 1)
		case driver.BUTTON_CALL_DOWN:
			fmt.Println("Order request DOWNWARD from floor:", msg.Floor+1, "\n")
			driver.Elev_set_button_lamp(msg.Button, msg.Floor, 1)
		case driver.BUTTON_COMMAND:
			fmt.Println("Order from inside the elevator to floor:", msg.Floor+1, "\n")
	}
	send_from_network <- msg
	fmt.Println("End of cost function \n")
}

func priorityHandler_getCost(client driver.Client) int {
	return 0
}

func priorityHandler_sort_all_ips() {
	
}