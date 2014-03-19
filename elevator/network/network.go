package network

import (
	driver "../driver"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var network_list [3][4]bool

func Network() {
	all_ips_m := make(map[string]time.Time)
	all_clients_m := make(map[string]driver.Client)

	msg_from_network := make(chan driver.Client)
	order_to_network := make(chan driver.Client, 1)
	order_from_network := make(chan driver.Client, 1)
	order_from_cost := make(chan driver.Client, 1)
	status_update_c := make(chan driver.Client, 1)
	send_lights_c := make(chan driver.Lights, 1)
	set_light_c := make(chan driver.Lights, 1)
	send_del_req_c := make(chan driver.Order, 1)
	del_order_c := make(chan driver.Order, 1)
	order_complete_c := make(chan driver.Order,1)

	localIP, _ := LocalIP()
	fmt.Println(localIP, "\n")

	init_elevator, init_hardware, current_floor := Initialize_elevator()
	if init_elevator && init_hardware {
		go driver.OrderHandler_process_orders(order_from_network, order_to_network, status_update_c, send_lights_c, send_del_req_c,order_complete_c ,current_floor, localIP)
	}

	go Read_msg(msg_from_network, set_light_c, del_order_c, localIP, all_clients_m)
	go Send_msg(order_to_network, send_lights_c, send_del_req_c)
	go Read_alive(all_ips_m, all_clients_m, localIP)
	go Send_alive(status_update_c)
	go Inter_process_communication(msg_from_network, order_from_network, order_from_cost, set_light_c, del_order_c, localIP, all_clients_m, order_complete_c)
	go Get_kill_sig()
	neverQuit := make(chan string)
	<-neverQuit
}

func Initialize_elevator() (init_elevator bool, init_hardware bool, prev driver.Order) {
	init_hardware = true
	if driver.Elev_init() == 0 {
		init_hardware = false
		fmt.Println("Unable to initialize elevator hardware\n")
	}
	fmt.Println("Press STOP button to stop elevator and exit program.\n")
	init_elevator, current_floor := driver.Elevator_init()
	if !init_elevator {
		fmt.Println("Unable to initialize elevator to floor\n")
	}
	return init_elevator, init_hardware, current_floor
}

func Inter_process_communication(msg_from_network chan driver.Client, order_from_network chan driver.Client, order_from_cost chan driver.Client, set_light_c chan driver.Lights, del_order_c chan driver.Order, localIP net.IP, all_clients map[string]driver.Client, order_complete_c chan driver.Order) {
	for {
		select {
		case new_order := <-msg_from_network:
			all_clients[new_order.Ip.String()] = new_order
			if new_order.Button != driver.BUTTON_COMMAND  {
				network_list[new_order.Button][new_order.Floor] = true
				priorityHandler(new_order, order_from_cost, all_clients)
			}
		case set_light := <-set_light_c:
			if set_light.Flag {
				driver.Elev_set_button_lamp(set_light.Button, set_light.Floor, 1)
			} else {
				driver.Elev_set_button_lamp(set_light.Button, set_light.Floor, 0)
			}
		case delete_order := <-del_order_c:
			network_list[delete_order.Button][delete_order.Floor] = false
			order_complete_c <- delete_order

		case send_order := <-order_from_cost:
			if send_order.Ip_from_cost.String() == localIP.String() {
				order_from_network <- send_order
			}
		case <-time.After(10 * time.Second):
			fmt.Println("timeout, 10 seconds has passed")
		}
	}
}

func Read_msg(msg_from_network chan driver.Client, set_light_c chan driver.Lights, del_order_c chan driver.Order, localIP net.IP, all_clients map[string]driver.Client) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20003")
	Check_error(err_conv_ip_listen)
	listener, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	var decoded_client driver.Client
	var decoded_lights driver.Lights
	var decoded_order driver.Order
	for {
		b := make([]byte, 1024)
		n, _, _ := listener.ReadFromUDP(b)
		//if raddr.IP.String() == localIP.String() {
		code := string(b[:3])
		switch code {
		case "cli":
			err_decoding := json.Unmarshal(b[3:n], &decoded_client)
			if err_decoding != nil {
				fmt.Println("error decoding client msg")
			}
			msg_from_network <- decoded_client

		case "lig":
			err_decoding := json.Unmarshal(b[3:n], &decoded_lights)
			if err_decoding == nil {
				set_light_c <- decoded_lights
			}
		case "del":
			err_decoding := json.Unmarshal(b[3:n], &decoded_order)
			if err_decoding == nil {
				del_order_c <- decoded_order
			}
		}
	}
}

func Send_msg(order_to_network chan driver.Client, send_lights_c chan driver.Lights, send_del_req_c chan driver.Order) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20003")
	Check_error(err_conv_ip)
	msg_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		select {
		case new_order := <-order_to_network:
			msg_encoded, err_encoding := json.Marshal(new_order)
			if err_encoding != nil {
				fmt.Println("error encoding json: ", err_encoding)
			}
			msg_encoded = append([]byte("cli"), msg_encoded...)
			msg_sender.Write(msg_encoded)
		case send_lights := <-send_lights_c:
			lights_encoded, err_encoding := json.Marshal(send_lights)
			if err_encoding != nil {
				fmt.Println("error encoding json: ", err_encoding)
			}
			lights_encoded = append([]byte("lig"), lights_encoded...)
			msg_sender.Write(lights_encoded)
		case send_del_req := <-send_del_req_c:
			delete_encoded, err_encoding := json.Marshal(send_del_req)
			if err_encoding != nil {
				fmt.Println("error encoding json: ", err_encoding)
			}
			delete_encoded = append([]byte("del"), delete_encoded...)
			msg_sender.Write(delete_encoded)
		}
	}
}

func Send_alive(status_update_c chan driver.Client) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	Check_error(err_conv_ip)
	alive_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		select {
		//time.Sleep(100*time.Millisecond)
		case status_update := <-status_update_c:
			status_encoded, err_encoding := json.Marshal(status_update)
			if err_encoding != nil {
				fmt.Println("error encoding json: ", err_encoding)
			}
			status_encoded = append([]byte("status"), status_encoded...)
			alive_sender.Write([]byte(status_encoded))
		case <-time.After(100 * time.Millisecond):
			alive_sender.Write([]byte("alive?"))
		}

	}
}

func Read_alive(all_ips map[string]time.Time, all_clients map[string]driver.Client, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	Check_error(err_conv_ip_listen)
	alive_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	var status_decoded driver.Client
	for {
		time.Sleep(50 * time.Millisecond)
		b := make([]byte, 1024)
		n, raddr, _ := alive_receiver.ReadFromUDP(b)
		code := string(b[:6])
		switch code {
		case "status":
			err_decoding := json.Unmarshal(b[6:n], &status_decoded)
			if err_decoding != nil {
				fmt.Println("error decoding client msg")
			}
			status_decoded.Ip = raddr.IP
			all_clients[raddr.IP.String()] = status_decoded
			Write_to_file(status_decoded)
		case "alive?":
			all_ips[raddr.IP.String()] = time.Now()
		}
		CheckForElapsedClients(all_ips, all_clients)
	}
}

func CheckForElapsedClients(all_ips map[string]time.Time, all_clients map[string]driver.Client) {
	for key, value := range all_ips {
		if time.Now().Sub(value) > 3*time.Second {
			fmt.Println("Deleting IP: ", key, " ", value)
			delete(all_ips, key)
			delete(all_clients, key)
		}
	}
}

func LocalIP() (net.IP, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}
			return v4, nil
		}
	}
	return nil, nil //errors.New("cannot find local IP address")
}
