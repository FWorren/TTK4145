package driver

import (
    "math"
)

const N_FLOORS = 4
const N_BUTTONS = 3 

type elev_button_type_t int

const  (
    BUTTON_CALL_UP elev_button_type_t = iota 
    BUTTON_CALL_DOWN
    BUTTON_COMMAND
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
    {FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
    {FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
    {FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4},
}

func Elev_init() int {
    if !Io_init() {
        return 0
    }

    for i := 0; i < N_FLOORS; i++ {
        if i != 0 {
            Elev_set_button_lamp(BUTTON_CALL_DOWN,i,0)
        }
        if i != N_FLOORS-1 {
            Elev_set_button_lamp(BUTTON_CALL_UP, i, 0)
        }
        Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
    }

    Elev_set_stop_lamp(0)
    Elev_set_door_open_lamp(0)
    Elev_set_floor_indicator(0)

    return 1
}

func Elev_set_speed(speed int){
    last_speed := 0
    if speed > 0 {
        Io_clear_bit(MOTORDIR)
    }else if speed < 0 {
        Io_set_bit(MOTORDIR)
    }else if last_speed < 0 {
        Io_clear_bit(MOTORDIR)
    }else if last_speed > 0 {
        Io_set_bit(MOTORDIR)
    }
    last_speed = speed
    Io_write_analog(MOTOR, int(2048+4*math.Abs(float64(speed))))
}

func Elev_get_floor_sensor_signal() int {
    if Io_read_bit(SENSOR1) {
        return 0
    }else if Io_read_bit(SENSOR2) {
        return 1
    }else if Io_read_bit(SENSOR3) {
        return 2
    }else if Io_read_bit(SENSOR4) {
        return 3
    }else {
        return -1
    }
}

func Elev_get_button_signal(button elev_button_type_t, floor int) int {
    // Need error handling before proceeding
    if Io_read_bit(button_channel_matrix[floor][int(button)]) {
        return 1
    }else {
        return 0
    }
}

func Elev_get_stop_signal() bool {
    return Io_read_bit(STOP)
}

func Elev_get_obstruction_signal() bool {
    return Io_read_bit(OBSTRUCTION)
}

func Elev_set_floor_indicator(floor int){
    // Need error handling before proceeding
    switch floor{
        case 0:
            Io_clear_bit(FLOOR_IND1)
            Io_clear_bit(FLOOR_IND2)
        case 1:
            Io_set_bit(FLOOR_IND1)
            Io_clear_bit(FLOOR_IND2)
        case 2:
            Io_clear_bit(FLOOR_IND1)
            Io_set_bit(FLOOR_IND2)
        case 3:
            Io_set_bit(FLOOR_IND1)
            Io_set_bit(FLOOR_IND2)
    }
}

func Elev_set_button_lamp(button elev_button_type_t, floor int, value int){
    // Need error handling before proceeding
    if value == 1 {
        Io_set_bit(lamp_channel_matrix[floor][int(button)]);
    }else {
        Io_clear_bit(lamp_channel_matrix[floor][int(button)]);        
    }
}

func Elev_set_stop_lamp(value int){
    if value == 1 {
        Io_set_bit(LIGHT_STOP);
    }else {
        Io_clear_bit(LIGHT_STOP);
    }
}

func Elev_set_door_open_lamp(value int){
    if value == 1 {
        Io_set_bit(DOOR_OPEN);
    }else {
        Io_clear_bit(DOOR_OPEN);
    }
}