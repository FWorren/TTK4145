package network

import (
	driver "../driver"
	"encoding/json"
	//"bufio"
	"log"
	"fmt"
	"io/ioutil"
	"os"
	"net"
	"strings"
)

func Write_to_file(client driver.Client) {
	client_encoded, err_encoding := json.Marshal(client)
	if err_encoding != nil {
		fmt.Println("error encoding json: ", err_encoding)
	}
	filename := "/tmp/data"
	filename += Get_last_ip_digits(client.Ip) + ".txt"
	err := ioutil.WriteFile(filename, client_encoded, 0644)
	if err != nil {
		log.Fatal(err)
	}	
}

func Read_file(file string) driver.Client {
	var decoded_client driver.Client
	file_opened, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1024)
	n,err1 := file_opened.Read(data)
	if err1 != nil {
		log.Fatal(err)
	}
	err_decoding := json.Unmarshal(data[:n], &decoded_client)
	if err_decoding != nil {
		fmt.Println("error decoding client from backup file")
	}
	return decoded_client
}

func Get_last_ip_digits(ip net.IP) string {
	ipstring := strings.Split(ip.String(),".")
	return ipstring[len(ipstring)-1]
}