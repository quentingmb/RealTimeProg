package main

import (
	"driver"
	"elevator"
	. "fmt"
	"extra"
	"network"
	"Elevatorlogic"
	"runtime"
	"time"
)

func main() {
	var myinfo network.Info
	var takerequest []network.Request

	runtime.GOMAXPROCS(runtime.NumCPU())

	myip := network.GetLocalIP()
	Println(myip)
	myinfo.Source = myip

	conf := extra.LoadConfig("./config/conf.json")

	generatedmessages_c := make(chan network.Networkmessage, 100)
	go network.TCPPeerToPeer(conf, myip, generatedmessages_c)

	state := "INIT"
	driver.IoInit()
	elevator.ElevInit()

	for {
		time.Sleep(10 * time.Millisecond)
		myinfo.State = state
		elevator.UpdateFloor()
		myinfo.LastFloor = elevator.CurrentFloor()
		network.NewInfo(myinfo, generatedmessages_c)
		switch state {
		case "INIT":
			{
				elevator.SetElevSpeed(-300)
			}
		case "IDLE":
			{
				elevator.SetElevSpeed(0)
			}
		case "UP":
			{
				elevator.SetElevSpeed(300)
			}
		case "DOWN":
			{
				elevator.SetElevSpeed(-300)
			}
		case "DOOR_OPEN":
			{
				elevator.SetElevDoorOpenLamp(1)
				for _, request := range takerequest {
					request.InOut = 0
					Println("Deleting request: ", request)
					time.Sleep(10 * time.Millisecond)
					network.Newrequest(generatedmessages_c, request)
				}
				elevator.SetElevSpeed(0)
				time.Sleep(3000 * time.Millisecond)
				elevator.etElevDoorOpenLamp(0)
			}
		case "ERROR":
			{
				elevator.SetElevSpeed(0)
			}
		}
		state, takerequest = ElevatorLogic.Nextstate(myip, conf.Elevators, myinfo.State)
	}
}
