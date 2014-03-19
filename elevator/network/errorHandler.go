package network

import (
	//driver "../driver"
	"fmt"
	//"bufio"
    //"io/ioutil"
    "os"
    "os/signal"
)

func Check_Connectivity() bool {
	return false
}

func Check_error(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func Get_kill_sig() {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt)
	<- sigchan
	fmt.Println("Program killed !")
	fmt.Println("write current status and orders to file!")
	os.Exit(0)
}
