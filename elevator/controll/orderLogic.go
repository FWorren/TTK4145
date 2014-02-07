package controll

var command [N_FLOORS]int

type Order struct {
	floor int
	dir int
}

var head_order Order
var previous_order Order

func OrderLogic_init() {
	for i := 0; i < N_FLOORS; i++ {
		command[i] = 0
	}
}

func OrderLogic_search_for_orders() {
	for i := 0; i < N_FLOORS; i++ {
		if driver.Elev_get_button_signal(BUTTON_COMMAND,i) {
			if command[i] != 1 {
				command[i] = 1
				driver.Elev_set_button_lamp(BUTTON_COMMAND,i,1)
			}
		}
	}
}

func OrderLogic_get_number_of_internal_orders() int {
	n_orders := 0
	for i := 0; i < N_FLOORS; i++ {
		if command[i] == 1 {
			n_orders += 1
		}
	}
	return n_orders
}

func OrderLogic_update_previous_order(floor int, direction int) {
	previous_order.floor = floor
	previous_order.dir = direction
}