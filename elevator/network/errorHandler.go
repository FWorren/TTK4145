package network

import (
	"fmt"
)

func Check_error(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}