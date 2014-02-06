package driver

import (
    "C"
)

type elev_button_type_t int
const  (
    BUTTON_CALL_UP elev_button_type_t = iota 
    BUTTON_CALL_DOWN 
    BUTTON_COMMAND 
)
 
func Elev_init() int {
    return int(C.elev_init())
}

func Elev_set_speed(speed int) int {
    C.elev_set_speed(C.int(speed))
}

func Elev_get_floor_sensor_signal() int {
    return int(C.elev_get_floor_sensor_signal())
}

func Elev_get_button_signal(button elev_button_type_t, floor int) int {
    return int(C.elev_get_button_signal(C.enum_button, C.int(floor)))
}

func Elev_get_stop_signal() int {
    return int(C.io_read_bit(STOP))
}

func Elev_get_obstruction_signal() int {
    return int(C.io_read_bit(OBSTRUCTION))
}

func Elev_set_floor_indicator(floor int){
    C.elev_set_floor_indicator(C.int(floor))
}

func Elev_set_button_lamp(button elev_button_type_t, floor int, value int){
    C.elev_set_button_lamp(C.enum_button, C.int(floor), C.int(value))
}

func Elev_set_stop_lamp(value int){
    C.elev_set_stop_lamp(C.int(value))
}

func Elev_set_door_open_lamp(value int){
    C.elev_set_door_open_lamp(C.int(value))
}