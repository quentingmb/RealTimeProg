package network

import (
	"drivers"
	"elevator"
	"encoding/json"
	"fmt"
	"log"
	"misc"
	"net"
	"os"
	"strings"
	"time"
	)

 type Request struct {
 	Direction elev.Elev_button
 	Floor int
 	Type int //internal or external
 	Ipsource string 
 }
 type Info struct {
 	PreviousFloor int
 	Floor int
 	State int // moving or not
 	IpSource string
 }
 type ElevatorMessage {
 	Request Request
 	Info   Info
 }

 type Con struct {
 	Address *net.TCPConn
	Connected bool
 }
var NoRequest=[]Request{Request{}}
var requestList=make([]Request,0)
var elevators = make(map[string]bool)
var infolist = make(map[string]Status)
var connections = make([]*net.TCPConn, 0)

func GetRequestList() []Request {
	return requestList
}

func UnpackElevatorMessage(packed []byte, error_c chan string) ElevatorMessage {
	var message ElevatorMessage
	err := json.Unmarshal(packed, &message)
	if err != nil {
		error_c <- "Error in unpacking the message: " + err.Error()
	}
	return message
}
//Check if there is any new orders, if it is it passes it to Neworder
func Requestdistr(generatedMsgs_c chan ElevatorMessage, myip string) {
	var button elevator.Elev_button
	for {
		for floor, buttons := range elevator.Button_channel_matrix {
			for butt, channel := range buttons {
				if drivers.ReadBit(channel) {
					if butt == 0 {
						button = elevator.BUTTON_CALL_UP
					} else if butt == 1 {
						button = elevator.BUTTON_CALL_DOWN
					} else {
						button = elevator.BUTTON_COMMAND
					}
					Neworder(generatedMsgs_c, Request{Direction: button, Floor: floor + 1, Type: 1, IpSourceource: myip})
					time.Sleep(time.Millisecond)
				}
			}
		}
	}
}


func Listener(conn *net.TCPListener, connect_c chan Con, error_c chan string) {
	for {
		newConn, err := conn.AcceptTCP()
		if err != nil {
			error_c <- "Accept trouble: " + err.Error()
		}
		connect_c <- Con{Address: newConn, Connect: true}
	}
}

func SendAliveMessages(connection *net.TCPConn, error_c chan string) {
	for {
		_, err := connection.Write([]byte("KEEPALIVE"))
		if err != nil {
			error_c <- "error in sending keepalive message: " + err.Error()
			return
		}
		time.Sleep(time.Second)
	}
}

func SendStatuslist(generatedMsgs_c chan ElevatorMessage) {
	myip := misc.GetLocalIP()
	myinfo := infolist[myip]
	generatedMsgs_c <- ElevatorMessage{Request: Request{}, Info: myinfo}
}

func Neworder(generatedMsgs_c chan ElevatorMessage, request Request) bool {
	if request.Direction != elevator.BUTTON_COMMAND {
		request.IpSource = ""
	}
	for _, r := range requestlist {
		if r == request {
			return false
		}
	}
	generatedMsgs_c <- ElevatorMessage{Request: request, Info: Info{}}
	return true
}

