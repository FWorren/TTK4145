package driver

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lcomedi -lm
// #include <png.h>
// #include "io.h"
import "C"

func Io_init() bool {
    return bool(int(C.io_init()) != -1)
}

func Io_set_bit(channel int) {
    C.io_set_bit(C.int(channel))
}

func Io_clear_bit(channel int) {
    C.io_clear_bit(C.int(channel))
}

func Io_write_analog(channel int, value int) {
    C.io_write_analog(C.int(channel),C.int(value))
}

func Io_read_bit(channel int) bool {
    return bool(int(C.io_read_bit(C.int(channel))) !=0 )
}

func Io_read_analog(channel int) int {
    return int(C.io_read_analog(C.int(channel)))
}