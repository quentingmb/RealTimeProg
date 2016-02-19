package main

import (
	"./driver"
	"fmt"
)

func main() {
	fmt.Println("Started")
	driver.Init(driver.ET_comedi)
	driver.Set_bit(driver.LIGHT_COMMAND1)
	
	/*for true {
        // Change direction when we reach top/bottom floor
        if driver.Get_floor_sensor_signal() == 4 - 1 {
            driver.Set_motor_direction(-1);
        } else if driver.Get_floor_sensor_signal() == 0 {
            driver.Set_motor_direction(1);
        }

        // Stop elevator and exit program if the stop button is pressed
        if driver.Get_stop_signal() ==1 {
            driver.Set_motor_direction(0);
            
        }
    }*/
driver.Set_motor_direction(0)

	fmt.Println("Done.")

	// We wait to make sure the driver starts all its threads & connections
	select {}

}
