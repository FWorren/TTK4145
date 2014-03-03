package network

import (
	"fmt"
	driver "../driver"
)

func cost(msg driver.Client, send_from_network chan driver.Client) {
	fmt.Println("Cost running \n")
	fmt.Println("Recieving order :", msg.Floor+1,"\n")
	Order_ext[msg.Button][msg.Floor] = true
	send_from_network <- msg
	fmt.Println("End of cost function \n")
}