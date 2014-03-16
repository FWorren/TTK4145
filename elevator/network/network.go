package network

import (
	"fmt"
	"net"
	//"strings"
	driver "../driver"
	"encoding/json"
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
	//set_lights_c := make(chan driver.Lights, 1)
	//set_light_c := make(chan driver.Lights, 1)
	del_orders_c := make(chan driver.Order, 1)
	del_order_c := make(chan driver.Order, 1)

	localIP, _ := LocalIP()
	fmt.Println(localIP, "\n")

	go Read_msg(msg_from_network, localIP, all_clients_m,set_light_c, del_order_c)
	go Send_msg(order_to_network, set_lights_c, del_orders_c)
	go Read_alive(all_ips_m, all_clients_m, localIP)
	go Send_alive()
	go Inter_process_communication(msg_from_network, order_from_network, order_from_cost, localIP, all_clients_m, del_order_c)

	init_elevator, init_hardware, current_floor := Initialize_elevator()
	if init_elevator && init_hardware {
		go driver.OrderHandler_process_orders(order_from_network, order_to_network, current_floor, localIP)
	}

	neverQuit := make(chan string)
	<-neverQuit
}

func Initialize_elevator() (init_elevator, init_hardware bool, prev driver.Order) {
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

func Inter_process_communication(msg_from_network chan driver.Client, order_from_network chan driver.Client, order_from_cost chan driver.Client, localIP net.IP, all_clients map[string]driver.Client, del_order_c chan driver.Order) {
	for {
		select {
		case new_order := <-msg_from_network:
			fmt.Println("msg_from_network: ", new_order.Ip.String())
			all_clients[new_order.Ip.String()] = new_order
			network_list[new_order.Button][new_order.Floor] = true
			priorityHandler(new_order, order_from_cost, all_clients)
		case delete_order := <-del_order_c:
			network_list[delete_order.Button][delete_order.Floor] = false
		case send_order := <-order_from_cost:
			driver.Elev_set_button_lamp(send_order.Button, send_order.Floor, 1)
			if send_order.Ip_from_cost.String() == localIP.String() {
				order_from_network <- send_order
				fmt.Println("order_from_network: ", send_order.Floor+1, "\n")
			}
		case <-time.After(10 * time.Second):
			fmt.Println("timeout, 10 seconds has passed")
		}
	}
}

func Read_msg(msg_from_network chan driver.Client, localIP net.IP, all_clients map[string]driver.Client, set_light_c chan driver.Lights, del_order_c chan driver.Order) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20003")
	Check_error(err_conv_ip_listen)
	listener, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	var decoded_client driver.Client
	var decoded_lights driver.Lights
	var decoded_deleted driver.Order
	for {
		b := make([]byte, 1024)
		n, raddr, _ := listener.ReadFromUDP(b)
		if raddr.IP.String() == localIP.String() { //HUSK DETTE HER !=
			err_decoding_client := json.Unmarshal(b[0:n], &decoded_client)
			err_decoding_lights := json.Unmarshal(b[0:n], &decoded_lights)
			err_decoding_deleted := json.Unmarshal(b[0:n], &decoded_deleted)
			if err_decoding_client == nil {
				fmt.Println("decoded: ",decoded_client)
				msg_from_network <- decoded_client
			}
			if err_decoding_lights == nil {
				set_light_c <- decoded_lights
			}
			if err_decoding_deleted == nil {
				del_order_c <- decoded_deleted
			}
		}
	}
}

func Send_msg(order_to_network chan driver.Client, set_lights_c chan driver.Lights, del_orders_c chan driver.Order) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20003")
	Check_error(err_conv_ip)
	msg_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		select {
		case new_order := <-order_to_network:
			fmt.Println("before json :",new_order)
			msg_encoded, err_encoding := json.Marshal(new_order)
			if err_encoding != nil {
				fmt.Println("error encoding json \n")
			}
			Check_error(err_encoding)
			msg_sender.Write(msg_encoded)
		case set_lights := <-set_lights_c:
			msg_encoded, err_encoding := json.Marshal(set_lights)
			if err_encoding != nil {
				fmt.Println("error encoding json \n")
			}
			Check_error(err_encoding)
			msg_sender.Write(msg_encoded)
		case del_orders := <-del_orders_c:
			msg_encoded, err_encoding := json.Marshal(del_orders)
			if err_encoding != nil {
				fmt.Println("error encoding json \n")
			}
			Check_error(err_encoding)
			msg_sender.Write(msg_encoded)
		}
	}
}

func Send_alive() {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	Check_error(err_conv_ip)
	alive_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		time.Sleep(1000 * time.Millisecond)
		alive_sender.Write([]byte("?"))
	}
}

func Read_alive(all_ips map[string]time.Time, all_clients map[string]driver.Client, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	Check_error(err_conv_ip_listen)
	alive_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	var client driver.Client
	for {
		time.Sleep(50 * time.Millisecond)
		b := make([]byte, 0)
		_, raddr, _ := alive_receiver.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			all_ips[raddr.IP.String()] = time.Now()
			client.Ip = raddr.IP
			all_clients[raddr.IP.String()] = client
			//fmt.Println("IP: ", raddr.IP.String(), " msg: ", string(b))
		} else {
			client.Ip = localIP
			all_clients[raddr.IP.String()] = client
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

func Check_error(err error) {
	if err != nil {
		fmt.Println("error:", err)
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

func Print_alive(all_ips map[string]string) {
	for key, value := range all_ips { // key = IP adress and value = time last seen
		fmt.Println("IP: ", key, " time: ", value, "\n")
	}
}
