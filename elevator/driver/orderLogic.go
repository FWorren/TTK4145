package driver

var command [N_FLOORS]int

type order_state_t int

const (
	UP order_state_t = iota
	DOWN
	SET
)

type Order struct {
	floor int
	dir int
}

var Head_order Order
var Previous_order Order

func OrderLogic_init() {
	for i := 0; i < N_FLOORS; i++ {
		command[i] = 0
	}
}

func OrderLogic_search_for_orders() {
	for i := 0; i < N_FLOORS; i++ {
		if Elev_get_button_signal(BUTTON_COMMAND,i) == 1 {
			if command[i] != 1 {
				command[i] = 1
				Elev_set_button_lamp(BUTTON_COMMAND,i,1)
			}
		}
	}
}

func OrderLogic_set_head_order() {
	//OrderLogic_get_direction()
	/*for {
		switch {
			case UP:

			case DOWN:
			default:

		}
	}*/
}

func OrderLogic_state_up() {
	
}

func OrderLogic_state_down() {
	
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
	Previous_order.floor = floor
	Previous_order.dir = direction
}