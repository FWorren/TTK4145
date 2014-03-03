package network

import (
	"fmt"
	"net"
	//"strings"
	"time"
	driver "../driver"
)

var Order_ext = [2][4]bool{
    {false, false, false, false},
    {false, false, false, false},
}

func Network() {
	msg_from_network := make(chan string)
	msg_to_network := make(chan string)
	order_to_network := make(chan driver.Client)
	order_from_network := make(chan driver.Client,10)
	send_from_network := make(chan driver.Client, 10)
	order_internal := make(chan driver.Client)
	//all_ips_m := make(map[string]time.Time)
	localIP, _ := LocalIP()
	fmt.Println(localIP)
	/*
	go Read_msg(msg_from_network, localIP)
	go Send_msg(msg_to_network)
	go Read_alive(all_ips_m, localIP)
	go Send_alive()
	*/

	go Inter_process_communication(msg_from_network,msg_to_network,order_to_network,order_from_network, send_from_network)
	Init_hardware(order_to_network, order_from_network, order_internal)
	
	neverQuit := make(chan string)
	<-neverQuit
}

func Init_hardware(order_to_network chan driver.Client, order_from_network chan driver.Client, order_internal chan driver.Client) {
	if driver.Elev_init() == 0 {
		fmt.Println("Unable to initialize elevator hardware\n")
	}

	fmt.Println("Press STOP button to stop elevator and exit program.\n")
	go driver.Elevator_statemachine()

	go driver.OrderHandler_process_orders(order_to_network, order_from_network, order_internal)
}

func Inter_process_communication(msg_from_network chan string, msg_to_network chan string, order_to_network chan driver.Client,order_from_network chan driver.Client, send_from_network chan driver.Client) {
	for {
		select {
			case <- msg_from_network:
				fmt.Println("msg_from_network")
			case <- msg_to_network:
				fmt.Println("msg_to_network")
			case msg :=<- order_to_network:
				fmt.Println("order_to_network")
				cost(msg, send_from_network)
				_ = msg
			case send_order := <- send_from_network:
				fmt.Println("order_from_network: ",send_order.Floor+1,"\n")
				fmt.Println(Order_ext, "\n")
				order_from_network <- send_order
			case <- time.After(5*time.Second):
				fmt.Println("timeout")
		}
	}
}

func Read_msg(msg_from_network chan driver.Client, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20003")
	Check_error(err_conv_ip_listen)
	listener, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	for {
		b := make([]byte, 1024)
		_, raddr, _ := listener.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			//msg_from_network <- strings.Trim(string(b), "\x00")
		}
	}
}

func Send_msg(msg_to_network chan string) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20003")
	Check_error(err_conv_ip)
	msg_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		msg_sender.Write([]byte(<-msg_to_network))
	}
}

func Send_alive() {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	Check_error(err_conv_ip)
	alive_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		time.Sleep(1000 * time.Millisecond)
		alive_sender.Write([]byte("Alive?"))
	}
}

func Read_alive(all_ips map[string]time.Time, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	Check_error(err_conv_ip_listen)
	alive_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	for {
		b := make([]byte, 1024)
		_, raddr, _ := alive_receiver.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			all_ips[raddr.IP.String()] = time.Now()
			for key, value := range all_ips {
				if time.Now().Sub(value) > 3*time.Second {
					delete(all_ips, key)
				}
			}
			fmt.Println("IP: ", raddr.IP.String(), " msg: ", string(b))
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