package network

import (
	"fmt"
)

func Check_Connectivity() bool {

}

func Check_error(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
