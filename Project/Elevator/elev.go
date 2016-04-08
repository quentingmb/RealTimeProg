//A completer 
//ajouter les imports

package elevator
import (
	"driver"
	
)

const N_FlOORS = 4
const N_BUTTONS = 4

type Elev_button int

var Lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{driver.LIGHT_UP1, driver.LIGHT_DOWN1, driver.LIGHT_COMMAND1},
	{driver.LIGHT_UP2, driver.LIGHT_DOWN2, driver.LIGHT_COMMAND2},
	{driver.LIGHT_UP3, driver.LIGHT_DOWN3, driver.LIGHT_COMMAND3},
	{driver.LIGHT_UP4, driver.LIGHT_DOWN4, driver.LIGHT_COMMAND4},
}

var Button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{driver.FLOOR_UP1, driver.FLOOR_DOWN1, driver.FLOOR_COMMAND1},
	{driver.FLOOR_UP2, driver.FLOOR_DOWN2, driver.FLOOR_COMMAND2},
	{driver.FLOOR_UP3, driver.FLOOR_DOWN3, driver.FLOOR_COMMAND3},
	{driver.FLOOR_UP4, driver.FLOOR_DOWN4, driver.FLOOR_COMMAND4},
}

func SetElevSpeed (speed int) {
	if speed == 0 {
		if driver.ReadBit(driver.MOTORDIR) {
			driver.ClearBit(driver.MOTORDIR)
		} else {
			driver.SetBit(driver.MOTORDIR)
		}
		time.Sleep(10 * time.Millisecond)
	}
	if speed > 0 {
		driver.ClearBit(driver.MOTORDIR)
	} else {
		driver.SetBit(driver.MOTORDIR)
	}
	driver.WriteAnalog(driver.MOTOR, 2048+4*math.Abs(speed))
}
func SetElevFloorIndicator(floor int) {
	//one light sould be on
	if floor == 3 || floor == 4 {
		driver.Set_bit(driver.FLOOR_IND1)
	} else {
		driver.Clear_bit(driver.FLOOR_IND1)
	}
	if floor == 2 || floor == 4 {
		driver.Set_bit(driver.FLOOR_IND2)
	} else {
		driver.Clear_bit(driver.FLOOR_IND2)
	}
}
func SetElevButtonLamp(button Elev_button, floor int, value int) {
	if value == 1 {
		driver.Set_bit(Lamp_channel_matrix[floor][button])
	} else {
		driver.Clear_bit(Lamp_channel_matrix[floor][button])
	}
}
func SetElevDoorOpenLamp(value int) {
	if value == 1 {
		driver.Set_bit(driver.DOOR_OPEN)
	} else {
		driver.Clear_bit(driver.DOOR_OPEN)
	}
}
func SetElevStopLamp(value int) {
	if value == 1 {
		driver.Set_bit(drivers.LIGHT_STOP)
	} else {
		driver.Clear_bit(drivers.LIGHT_STOP)
	}
}
func GetElevFloorSensorSignal() int {
	if driver.Read_bit(driver.SENSOR1) {
		return 1
	} else if driver.Read_bit(driver.SENSOR2) {
		return 2
	} else if driver.Read_bit(driver.SENSOR3) {
		return 3
	} else if driver.Read_bit(driver.SENSOR4) {
		return 4
	} else {
		return -1
	}

}

func GetElevButtonSignal(button Elev_button, floor int) int {
	if driver.Read_bit(Button_channel_matrix[floor][button]) {
		return 1
	} else {
		return 0
	}
}
func GetElevStopSignal() bool {
	return driver.Read_bit(driver.STOP)
}
func GetElevObstructionSignal() bool {
		return driver.Read_bit(driver.OBSTRUCTION)
}
func ElevInit() {
	// Zero all floor button lamps
	for i := 1; i < N_FLOORS; i++ {
		for b:=0;b<N_BUTTONS;b++ {
			SetElevButtonLamp(b, i, 0)
		}		
	}

	// Clear stop lamp, door open lamp, and set floor indicator to ground floor.
	ElevSetStopLamp(0)
	ElevSetDoorOpenLamp(0)
	ElevSetFloorIndicator(0)

}
