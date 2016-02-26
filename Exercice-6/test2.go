package main

import (
	."fmt"
	"net"
	"time"
)

func count(counter int) {
	var data = make([]byte, 256)
	BROADCAST_IPv4 := net.IPv4(129, 241, 187, 255)
	port := 58017
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error SendAliveMessage 1")
	}
	
	for {
		counter++
		Println(counter)
		
		data[0] = byte(counter)
		_, err := socket.Write(data)
		if err != nil {
			Printf("error SendAliveMessage 2")
		}
		time.Sleep(1000*time.Millisecond)
	}
}

func SendAliveMessage() {
	BROADCAST_IPv4 := net.IPv4(129, 241, 187, 255)
	port := 57017
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error SendAliveMessage 1")
	}
	
	for {
		
		data := []byte(GetLocalIP())
		_, err := socket.Write(data)
		if err != nil {
			Printf("error SendAliveMessage 2")
		}
		time.Sleep(100*time.Millisecond)
	}
}

func newPrimary(counter int){
	time.Sleep(time.Second)
	go SendAliveMessage()
	go count(counter)
}


func Receive(startChan chan string){
	addr, _ := net.ResolveUDPAddr("udp4", ":58017")
	socket, err := net.ListenUDP("udp4", addr)
	if err != nil {
		Printf("error ReceiveMessage 1")
	}
	var counter int = 0
	data := make([]byte, 256)
	
	for {
		select {
		case msg := <- startChan:
			msg = msg
			socket.Close()
			go newPrimary(counter)
			return
			
		default:
		  
		    
			tall := socket.SetReadDeadline(time.Now().Add(3*time.Second)) //break if no message in 2.4 seconds
			tall = tall
			_,_,err := socket.ReadFromUDP(data)
				
			if err != nil {
				Printf("error ReceiveMessage 2\n")
			}
			counter = int(data[0])
			Println(counter)
		}
	}
}

func ReceiveAliveMessage(recChan chan string, startChan chan string){
	addr, _ := net.ResolveUDPAddr("udp4", ":57017")
	socket, err := net.ListenUDP("udp4", addr)
	if err != nil {
		Printf("error ReceiveAliveMessage 1\n")
	}
	data := make([]byte, 256)
	
	for {
		
		errorrr := socket.SetReadDeadline(time.Now().Add(600*time.Millisecond))
		if errorrr != nil {
			socket.Close()
			return
		}
		_,_,err := socket.ReadFromUDP(data)
		
		
		
		if err != nil {
			Printf("error ReceiveAliveMessage 2\n")
		}
		recChan <- "alive"
		if string(data[:15]) == GetLocalIP() {
			//println("msg received")

		}
		println("msg received")

	}
}

func Timeout(recChan chan string, startChan chan string){
	for{
		select{
		case <- recChan:
			//Println("primary alive")
			
		case <-time.After(500*time.Millisecond):
			Println("hei")
			startChan <- "start"
			
			return
		}
	}
}



func GetLocalIP() (localIP string) {
	addrs, err := net.InterfaceAddrs()
    if err != nil {
    	Println(err)
    }
    
    for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
		    	localIP = ipnet.IP.String()
			}
		}
    }
    
	return
}

func main() {
	recChan := make(chan string)
	startChan := make(chan string)
	doneChan := make(chan string)
	
	
	go ReceiveAliveMessage(recChan, startChan)
	go Receive(startChan)
	time.Sleep(100*time.Millisecond)
	go Timeout(recChan, startChan)
	println(<-doneChan)
}

