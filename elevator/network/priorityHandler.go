package network

import (
	driver "../driver"
	"fmt"
)

func getCost(msg driver.Client, send_from_network chan driver.Client) {
	fmt.Println("  Cost running \n")
	switch msg.Button {
	case driver.BUTTON_CALL_UP:
		fmt.Println("Order request UPWARD from floor:", msg.Floor+1, "\n")
	case driver.BUTTON_CALL_DOWN:
		fmt.Println("Order request DOWNWARD from floor:", msg.Floor+1, "\n")
	case driver.BUTTON_COMMAND:
		fmt.Println("Order from inside the elevator to floor:", msg.Floor+1, "\n")
	}
	//fmt.Println("Order_list : \n",msg.Order_list,"\n")
	//msg.Order_list[msg.Button][msg.Floor] = true
	fmt.Println("This is where the cost function will")
	fmt.Println("decide which elevator will handle the order\n")
	send_from_network <- msg
	fmt.Println("End of cost function \n")
}
