package network

import (
	"net"
	"fmt"
	"time"
	"strings"
)

func network() {
	msg_from_network := make(chan string)
	msg_to_network := make(chan string)
	all_ips_m := make(map[string]string)
	initialization(msg_from_network,msg_to_network,all_ips_m)
	
	neverQuit := make(chan string)
	<-neverQuit
}

func initialization(msg_from_network chan string, msg_to_network chan string, all_ips map[string]string) {
	localIP,_ := localIP()
	fmt.Println(localIP)
	go read_msg(msg_from_network, localIP)
	go send_msg(msg_to_network)
	go read_alive(all_ips, localIP)
	go send_alive()
}

func read_msg(msg_from_network chan string, localIP net.IP){
		laddr , err_conv_ip_listen := net.ResolveUDPAddr("udp",":20003")
		check_error(err_conv_ip_listen)
		listener, err_listen := net.ListenUDP("udp", laddr)
		check_error(err_listen)
		
		for {
			b := make([]byte, 1024)
			_,raddr,_ := listener.ReadFromUDP(b)
			if raddr.IP.String() != localIP.String() {
				msg_from_network <- strings.Trim(string(b), "\x00")
			}
		}
}

func send_msg(msg_to_network chan string) {
	baddr,err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20003")
	check_error(err_conv_ip)
	msg_sender, err_dialudp := net.DialUDP("udp",nil,baddr)
	check_error(err_dialudp)
	for {
		msg_sender.Write([]byte(<-msg_to_network))
	}
}

func send_alive() {
	baddr,err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	check_error(err_conv_ip)
	alive_sender, err_dialudp := net.DialUDP("udp",nil,baddr)
	check_error(err_dialudp)
	for {
		time.Sleep(1000*time.Millisecond)
		alive_sender.Write([]byte("Alive?"))
	}
}

func read_alive(all_ips map[string]string, localIP net.IP) {
		laddr , err_conv_ip_listen := net.ResolveUDPAddr("udp",":20060")
		check_error(err_conv_ip_listen)
		alive_receiver, err_listen := net.ListenUDP("udp", laddr)
		check_error(err_listen)
		for {
			b := make([]byte, 1024)
			_,raddr,_ := alive_receiver.ReadFromUDP(b)
			if raddr.IP.String() != localIP.String() {
				all_ips[raddr.IP.String()] = time.Now().String()
				fmt.Println("IP: ", raddr.IP.String(), " msg: ", string(b))
			}
		}
}

func check_error(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())	
	}
}

func localIP() (net.IP, error) {
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

func print_alive(all_ips map[string]string) {
	for key,value := range all_ips { // key = IP adress and value = time last seen
		fmt.Println("IPaddress: ", key, " time: ", value, "\n")
	}
}