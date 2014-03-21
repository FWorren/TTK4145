package network

import (
	driver "../driver"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func Restore_floorpanel_orders(ip net.IP) (bool, driver.Client) {
	ip_last_digits := Get_last_ip_digits(ip)
	file := "network/backup/data"
	file += ip_last_digits
	err, backup_client := Read_file(file)
	fmt.Println("reading")
	if !err {
		return true, backup_client
	}
	fmt.Println("Error reading from file")
	return false, backup_client
}

func Search_for_lost_orders(client driver.Client, order_from_cost chan driver.Client, all_clients map[string]driver.Client) {
	for i := 0; i < driver.N_FLOORS; i++ {
		if client.Order_list[driver.BUTTON_CALL_DOWN][i] {
			client.Floor = i
			client.Button = driver.BUTTON_CALL_DOWN
			priorityHandler(client, order_from_cost, all_clients)
		}
		if client.Order_list[driver.BUTTON_CALL_UP][i] {
			client.Floor = i
			client.Button = driver.BUTTON_CALL_UP
			priorityHandler(client, order_from_cost, all_clients)
		}
	}
}

func Restore_command_orders(check_backup_c chan driver.Client, localIP net.IP) bool {
	ip_last_digits := Get_last_ip_digits(localIP)
	file := "network/backup/data"
	file += ip_last_digits
	err, backup_client := Read_file(file)
	if !err {
		check_backup_c <- backup_client
		return true
	}
	return false
}

func Sync_lights(all_clients map[string]driver.Client, localIP net.IP) {
	for key, value := range all_clients {
		if key != localIP.String() {
			for i := 0; i < driver.N_FLOORS; i++ {
				if value.Order_list[driver.BUTTON_CALL_UP][i] {
					driver.Elev_set_button_lamp(driver.BUTTON_CALL_UP, i, 1)
				}
				if value.Order_list[driver.BUTTON_CALL_DOWN][i] {
					driver.Elev_set_button_lamp(driver.BUTTON_CALL_DOWN, i, 1)
				}
			}
			break
		}
	}
}

func Check_connectivity(disconnected chan int, netstate_c chan driver.NetState_t) {
	connected := make(chan bool)

	go func() {
		timeOut := make(<-chan time.Time)
		for {
			timeOut = time.After(1 * time.Second)
			select {
			case <-connected:
				netstate_c <- 1
				break
			case <-timeOut:
				netstate_c <- 0
				disconnected <- 1
			}
		}
	}()

	for {
		time.Sleep(200 * time.Millisecond)
		_, err := net.LookupHost("google.com")
		if err == nil {
			connected <- true
		}
	}
}

func Get_kill_sig() {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	fmt.Println("Program killed !")
	driver.Elev_set_speed(0)
	os.Exit(0)
}

func Check_error(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
