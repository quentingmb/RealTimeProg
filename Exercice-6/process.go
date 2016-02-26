package main

import (
	"fmt"
	"net"
	"time"
	"os/exec"
	"strconv"
	"strings"
)

const (
	host = "129.241.187.255"
	port = "20013"
)

func SpawnProcess() {
	fmt.Println("Spawning backup")
	cmd := exec.Command("gnome-terminal", "-x", "go", "run", "process.go")
	out, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return
	}
	print(string(out))
}

func backup (sock *net.UDPConn, primaryAlive bool, counter *int) bool{
	for(primaryAlive){
		sock.SetReadDeadline(time.Now().Add(2*time.Second))
		data := make([]byte, 255)
		n, _, err := sock.ReadFromUDP(data[0:])
		if err != nil {
			primaryAlive = false
			return primaryAlive
		} else {
			SpawnProcess()
			s := strings.TrimLeft(string(data[:n]), "Count: ")
			count,_ := strconv.Atoi(s)
			*counter = count
			
			fmt.Println("Backup, count:", *counter)
		}
	}
	return true
}

func main() {
	primaryAlive := false
	counter := 0
	t_count := 0
	udpAddr, _ := net.ResolveUDPAddr("udp", host + ":" + port)
	sock, _ := net.ListenUDP("udp", udpAddr)
	primaryAlive = backup(sock, primaryAlive, &counter)
	sock.Close()
    
	t_count = counter
	fmt.Println(t_count)
	addr, _ := net.ResolveUDPAddr("udp4", host + ":" + port)
	sock2, _ := net.DialUDP("udp4", nil, addr)

	for {
		msg := "Count:" + strconv.Itoa(counter)
		_, err := sock2.Write([]byte(msg))
		fmt.Println("Primary count: ", counter)
		if err != nil {
			fmt.Println("Error:Broadcast", err.Error())
		}
		if (counter == 5 + t_count) {
			t_count=5
			counter=t_count+1
			break
		}
		counter++	
		time.Sleep(time.Second)
	}
	sock2.Close()
	SpawnProcess()
}
