package network

import (
	"fmt"
	"net"
	//"strings"
	driver "../driver"
	"encoding/json"
	"time"
)

func Network() {
	msg_from_network := make(chan driver.Client)
	order_to_network := make(chan driver.Client, 10)
	order_from_network := make(chan driver.Client, 10)
	order_from_cost := make(chan driver.Client, 10)
	order_internal := make(chan driver.Client)
	all_ips_m := make(map[string]time.Time)
	all_clients_m := make(map[string]driver.Client)
	localIP, _ := LocalIP()
	fmt.Println(localIP, "\n")
	go Read_msg(msg_from_network, localIP, all_clients_m)
	go Send_msg(order_to_network)
	go Read_alive(all_ips_m, localIP)
	go Send_alive()
	go Inter_process_communication(msg_from_network, order_from_network, order_from_cost, localIP, all_clients_m)
	Init_hardware(order_from_network, order_to_network, order_internal, localIP)

	neverQuit := make(chan string)
	<-neverQuit
}

func Init_hardware(order_from_network chan driver.Client, order_to_network chan driver.Client, order_internal chan driver.Client, localIP net.IP) {
	if driver.Elev_init() == 0 {
		fmt.Println("Unable to initialize elevator hardware\n")
	}
	fmt.Println("Press STOP button to stop elevator and exit program.\n")
	go driver.Elevator_statemachine()
	var new_client driver.Client
	driver.Init_orderlist(new_client)
	go driver.OrderHandler_process_orders(order_from_network, order_to_network, order_internal, localIP)
}

func Inter_process_communication(msg_from_network chan driver.Client, order_from_network chan driver.Client, order_from_cost chan driver.Client, localIP net.IP, all_clients map[string]driver.Client) {
	for {
		select {
		case new_order := <-msg_from_network:
			fmt.Println("msg_from_network: ", new_order.Ip.String())
			priorityHandler(new_order, order_from_cost, all_clients)
		case send_order := <-order_from_cost:
			if send_order.Ip.String() == localIP.String() {
				order_from_network <- send_order
				fmt.Println("order_from_network: ", send_order.Floor+1, "\n")
			}
		case <-time.After(10 * time.Second):
			fmt.Println("timeout, 10 seconds has passed")
		}
	}
}

func Read_msg(msg_from_network chan driver.Client, localIP net.IP, all_clients map[string]driver.Client) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20003")
	Check_error(err_conv_ip_listen)
	listener, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	var msg_decoded driver.Client
	for {
		b := make([]byte, 1024)
		n, raddr, _ := listener.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			err_decoding := json.Unmarshal(b[0:n], &msg_decoded)
			if err_decoding != nil {
				fmt.Println("error DECODING order \n")
			}
			Check_error(err_decoding)
			fmt.Println("decoded msg: ", msg_decoded)
			all_clients[msg_decoded.Ip.String()] = msg_decoded
			msg_from_network <- msg_decoded
		}
	}
}

func Send_msg(order_to_network chan driver.Client) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20003")
	Check_error(err_conv_ip)
	msg_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		select {
		case new_order := <-order_to_network:
			msg_encoded, err_encoding := json.Marshal(new_order)
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

func Read_alive(all_ips map[string]time.Time, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	Check_error(err_conv_ip_listen)
	alive_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	for {
		time.Sleep(50 * time.Millisecond)
		b := make([]byte, 0)
		_, raddr, _ := alive_receiver.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			all_ips[raddr.IP.String()] = time.Now()
			fmt.Println("IP: ", raddr.IP.String(), " msg: ", string(b))
		}
		CheckForElapsedClients(all_ips)
	}
}

func CheckForElapsedClients(all_ips map[string]time.Time) {
	for key, value := range all_ips {
		if time.Now().Sub(value) > 3*time.Second {
			fmt.Println("Deleting IP: ", key, " ", value)
			delete(all_ips, key)
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
