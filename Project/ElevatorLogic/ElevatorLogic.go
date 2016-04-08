
package ElevatorLogic

import (
	"Elevator"
	"extra"
	"Network"
)
//This function does a BFS-search through all orders to find the most effective solution
func Nextrequest(myip string, Elevatorlist []misc.Elevator) Network.Request {
	var statelist = make(map[string]network.Info)
	infolist := Network.GetInfoList()
	for host, info := range infolist {
		statelist[host] = info
	}
	requestlist := Network.GetRequestList()
insideloop:
	for _, request := range requestlist {
		if request.Direction != Elevator.BUTTON_COMMAND {
			continue insideloop
		}
		for _, elevator := range Elevatorlist {
			if info, ok := statelist[Elevator.Address]; ok {
				if ((info.State == "UP" || info.State == "IDLE") && info.LastFloor <= request.Floor) || ((info.State == "DOWN" || info.State == "IDLE") && info.LastFloor >= request.Floor) {
					if info.Source == request.Source {
						if info.Source == myip {
							return request
						} else {
							delete(statelist, Elevator.Address)
							continue insideloop
						}
					}
				}
			}
		}
		for _, elevator := range Elevatorlist {
			if info, ok := statelist[Elevator.Address]; ok {
				if (info.State == "UP" && info.LastFloor >= request.Floor) || (info.State == "DOWN" && info.LastFloor <= request.Floor){
					if info.Source == request.Source {
						if info.Source == myip {
							return request
						} else {
							delete(statelist, Elevator.Address)
							continue insideloop
						}
					}
				}
			}
		}
	}
requestloop:
	for _, request := range requestlist {
		if request.Direction == elevator.BUTTON_COMMAND {
			continue requestloop
		}
		for i := 0; i < elevator.N_FLOORS; i++ {
			for _, elevator := range Elevatorlist {
				if info, ok := statelist[Elevator.Address]; ok {
					if i != 0 && (info.State == "UP" && info.LastFloor+i == request.Floor) || (info.State == "DOWN" && info.LastFloor-i == request.Floor) {
						if statelist[Elevator.Address].Source == myip {
							return request
						} else {
							delete(statelist, Elevator.Address)
							continue requestloop
						}
					}
				}
			}
			for _, elevator := range Elevatorlist {
				if info, ok := statelist[Elevator.Address]; ok {
					if info.State == "IDLE" && (info.LastFloor == request.Floor+i || info.LastFloor == request.Floor-i) {
						if statelist[Elevator.Address].Source == myip {
							return request
						} else {
							delete(statelist, Elevator.Address)
							continue requestloop
						}
					}
				}
			}
		}
	}
	return network.EmptyRequest[0]
}

//This function return orders the elevator should stop for
func Stop(myip string, mystate string) []Network.Request {
	var takerequest []Network.Request
	requestlist := Network.GetRequestList()
	for _, request := range requestlist {
		if (request.Direction == Elevator.BUTTON_COMMAND && request.Source == myip) || (request.Direction == Elevator.BUTTON_CALL_UP && mystate == "UP") || (request.Direction == Elevator.BUTTON_CALL_DOWN && mystate == "DOWN") {
			if request.Floor == Elevator.CurrentFloor() && elevator.AtFloor() {
				takerequest = append(takerequest, request)
			}
		}
	}
	return takerequest
}
//This function returns the next state for the elevator
func Nextstate(myip string, elevators []extra.Elevator, mystate string) (string, []Network.Request) {
	if Elevator.GetElevObstructionSignal() {
		Eelevator.SetElevStopLamp(1)
		return "ERROR", nil
	} else if mystate == "ERROR" {
		Elevator.SetElevStopLamp(0)
		return "INIT", nil
	}

	stop := Stop(myip, mystate)
	if len(stop) != 0 {
		return "DOOR_OPEN", stop
	}

	next := Nextrequest(myip, elevators)
	if elevator.AtFloor() && next.Floor == Elevator.CurrentFloor() {
		return "DOOR_OPEN", append(stop, next)
	}
	if next.Floor > Elevator.CurrentFloor() {
		return "UP", nil
	} else if next.Floor < Elevator.CurrentFloor() && next.Floor != 0 {
		return "DOWN", nil
	} else if Elevator.AtFloor() {
		return "IDLE", nil
	} else {
		return mystate, nil
	}
}
