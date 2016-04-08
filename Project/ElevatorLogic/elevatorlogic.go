
package elevatorlogic

import (
	"elevator"
	"config"
	"network"
)
//This function does a BFS-search through all orders to find the most effective solution
func Nextrequest(myip string, Elevatorlist []misc.Elevator) network.Request {
	var statelist = make(map[string]network.Info)
	infolist := network.GetInfoList()
	for host, info := range infolist {
		statelist[host] = info
	}
	requestlist := network.GetRequestList()
insideloop:
	for _, request := range requestlist {
		if request.Direction != elevator.BUTTON_COMMAND {
			continue insideloop
		}
		for _, elevator := range Elevatorlist {
			if info, ok := statelist[elevator.Address]; ok {
				if ((info.State == "UP" || info.State == "IDLE") && info.LastFloor <= request.Floor) || ((info.State == "DOWN" || info.State == "IDLE") && info.LastFloor >= request.Floor) {
					if info.Source == request.Source {
						if info.Source == myip {
							return request
						} else {
							delete(statelist, elevator.Address)
							continue insideloop
						}
					}
				}
			}
		}
		for _, elevator := range Elevatorlist {
			if info, ok := statelist[elevator.Address]; ok {
				if (info.State == "UP" && info.LastFloor >= request.Floor) || (info.State == "DOWN" && info.LastFloor <= request.Floor){
					if info.Source == request.Source {
						if info.Source == myip {
							return request
						} else {
							delete(statelist, elevator.Address)
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
				if info, ok := statelist[elevator.Address]; ok {
					if i != 0 && (info.State == "UP" && info.LastFloor+i == request.Floor) || (info.State == "DOWN" && info.LastFloor-i == request.Floor) {
						if statelist[elevator.Address].Source == myip {
							return request
						} else {
							delete(statelist, elevator.Address)
							continue requestloop
						}
					}
				}
			}
			for _, elevator := range Elevatorlist {
				if info, ok := statelist[elevator.Address]; ok {
					if info.State == "IDLE" && (info.LastFloor == request.Floor+i || info.LastFloor == request.Floor-i) {
						if statelist[elevator.Address].Source == myip {
							return request
						} else {
							delete(statelist, elevator.Address)
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
func Stop(myip string, mystate string) []network.Request {
	var takerequest []network.Request
	requestlist := network.GetRequestList()
	for _, request := range requestlist {
		if (request.Direction == elevator.BUTTON_COMMAND && request.Source == myip) || (request.Direction == elevator.BUTTON_CALL_UP && mystate == "UP") || (request.Direction == elevator.BUTTON_CALL_DOWN && mystate == "DOWN") {
			if request.Floor == elevator.CurrentFloor() && elevator.AtFloor() {
				takerequest = append(takerequest, request)
			}
		}
	}
	return takerequest
}
//This function returns the next state for the elevator
func Nextstate(myip string, elevators []misc.Elevator, mystate string) (string, []network.Request) {
	if elevator.GetElevObstructionSignal() {
		elevator.SetElevStopLamp(1)
		return "ERROR", nil
	} else if mystate == "ERROR" {
		elevator.SetElevStopLamp(0)
		return "INIT", nil
	}

	stop := Stop(myip, mystate)
	if len(stop) != 0 {
		return "DOOR_OPEN", stop
	}

	next := Nextrequest(myip, elevators)
	if elevator.AtFloor() && next.Floor == elevator.CurrentFloor() {
		return "DOOR_OPEN", append(stop, next)
	}
	if next.Floor > elevator.CurrentFloor() {
		return "UP", nil
	} else if next.Floor < elevator.CurrentFloor() && next.Floor != 0 {
		return "DOWN", nil
	} else if elevator.AtFloor() {
		return "IDLE", nil
	} else {
		return mystate, nil
	}
}
