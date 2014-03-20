package network

import (
	driver "../driver"
	"encoding/json"
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

	filename := "network/backup/data"
	filename += Get_last_ip_digits(client.Ip)

	_, err := os.Open(filename)
	if(err != nil){
		fmt.Println("error opening file")
		_, _ = os.Create(filename)
	}

	err = ioutil.WriteFile(filename, client_encoded, 0644 )
	if err != nil {
		fmt.Println("error writing to file")
		log.Fatal(err)
	}
}

func Read_file(file string) (bool,driver.Client) {
	var decoded_client driver.Client

	file_opened, err := os.Open(file)
	if(err != nil){
		fmt.Println("error opening file")
		_, _ = os.Create(file)
		return true, decoded_client
	}
	data := make([]byte, 1024)
	n,err1 := file_opened.Read(data)
	if err1 != nil {
		fmt.Println("error reading file")
		fmt.Println(err1)

	}

	err_decoding := json.Unmarshal(data[:n], &decoded_client)
	if err_decoding != nil {
		fmt.Println("error decoding client from backup file")
	}
	return false,decoded_client
}

func Get_last_ip_digits(ip net.IP) string {
	ipstring := strings.Split(ip.String(),".")
	return ipstring[len(ipstring)-1]
}