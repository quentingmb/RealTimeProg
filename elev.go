package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
#include "channels.h"
*/
import "C"

func Set_motor_direction(dirn int) {
	 C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}


func Set_button_lamp(button int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor));
}

func Set_door_open_lamp(value int)  {
	C.elev_set_door_open_lamp(C.int(value))
}

func Set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func Get_button_signal(button int , floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}
func Get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}
func Get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}
func Get_obstruction_signal() int {
	return int(C.elev_get_obstruction_signal())
}
