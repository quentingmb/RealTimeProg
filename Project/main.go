package main

import (
	"driver"
	"Elevator"
	. "fmt"
	"extra"
	"Network"
	"ElevatorLogic"
	"runtime"
	"time"
)

func main() {
	var myinfo Network.Info
	var takerequest []Network.Request

	runtime.GOMAXPROCS(runtime.NumCPU())

	myip := Network.GetLocalIP()
	Println(myip)
	myinfo.Source = myip

	conf := extra.LoadConfig("./config/conf.json")

	generatedmessages_c := make(chan Network.Networkmessage, 100)
	go Network.TCPPeerToPeer(conf, myip, generatedmessages_c)

	state := "INIT"
	driver.IoInit()
	Elevator.ElevInit()

	for {
		time.Sleep(10 * time.Millisecond)
		myinfo.State = state
		Elevator.UpdateFloor()
		myinfo.LastFloor = Elevator.CurrentFloor()
		Network.NewInfo(myinfo, generatedmessages_c)
		switch state {
		case "INIT":
			{
				Elevator.SetElevSpeed(-300)
			}
		case "IDLE":
			{
				Elevator.SetElevSpeed(0)
			}
		case "UP":
			{
				Elevator.SetElevSpeed(300)
			}
		case "DOWN":
			{
				Elevator.SetElevSpeed(-300)
			}
		case "DOOR_OPEN":
			{
				Elevator.SetElevDoorOpenLamp(1)
				for _, request := range takerequest {
					request.InOut = 0
					Println("Deleting request: ", request)
					time.Sleep(10 * time.Millisecond)
					Network.Newrequest(generatedmessages_c, request)
				}
				Elevator.SetElevSpeed(0)
				time.Sleep(3000 * time.Millisecond)
				Elevator.etElevDoorOpenLamp(0)
			}
		case "ERROR":
			{
				Elevator.SetElevSpeed(0)
			}
		}
		state, takerequest = ElevatorLogic.Nextstate(myip, conf.Elevators, myinfo.State)
	}
}
