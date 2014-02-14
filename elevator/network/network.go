package network

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func Network() {
	msg_from_network := make(chan string)
	msg_to_network := make(chan string)
	all_ips_m := make(map[string]time.Time)
	Initialization(msg_from_network, msg_to_network, all_ips_m)

	/*neverQuit := make(chan string)
	<-neverQuit*/
}

func Initialization(msg_from_network chan string, msg_to_network chan string, all_ips map[string]time.Time) {
	localIP, _ := LocalIP()
	fmt.Println(localIP)
	go Read_msg(msg_from_network, localIP)
	go Send_msg(msg_to_network)
	go Read_alive(all_ips, localIP)
	go Send_alive()
}

func Read_msg(msg_from_network chan string, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20003")
	Check_error(err_conv_ip_listen)
	listener, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)

	for {
		b := make([]byte, 1024)
		_, raddr, _ := listener.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			msg_from_network <- strings.Trim(string(b), "\x00")
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
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20060")
	Check_error(err_conv_ip_listen)
	alive_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	for {
		b := make([]byte, 1024)
		_, raddr, _ := alive_receiver.ReadFromUDP(b)
		if raddr.IP.String() != localIP.String() {
			all_ips[raddr.IP.String()] = time.Now()

			for key, value := range all_ips {
				err_dead_conn := alive_receiver.SetReadDeadline(value.Add(3 * time.Second))
				if err_dead_conn != nil {
					delete(all_ips, key)
				}
			}
			fmt.Println("IP: ", raddr.IP.String(), " msg: ", string(b))
		}
	}
}

func Check_error(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())
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
		fmt.Println("IPaddress: ", key, " time: ", value, "\n")
	}
}
