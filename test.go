package main

import(
	"fmt"
	"os"
	"os/exec"
	"net"
	"strconv"
	"strings"
	"time"
)

var Addr string = "129.241.187.255"
var port string = ":20013"

func checkErr(err error){
	if err != nil {
		fmt.Println("An unrecovarable error occured", err.Error())
		os.Exit(0)
	}
}

func spawnProcess(){
	cmd := exec.Command("gnome-terminal", "-x", "go", "run", "process.go")
	out, err := cmd.Output()
	checkErr(err)
	fmt.Println(string(out))
}

func getCount(msg string) int{
	n := strings.TrimLeft(msg, "Count: ")
	count, err := strconv.Atoi(n)
	if(err != nil){
		return -1
	}
	return count
}

func backupProcess(conn *net.UDPConn, primaryAlive bool, count *int) bool{
	for(primaryAlive){
		conn.SetReadDeadline(time.Now().Add(time.Second*2)) //takes some time to open terminal
		data := make([]byte, 16)
		length, _, err := conn.ReadFromUDP(data[0:])
		if err != nil {
			primaryAlive = false
			return primaryAlive
		} else{
			*count = getCount(string(data[0:length]))
			fmt.Println("Backup,  count:", *count)
		}
	}
	return true
}


func main(){

	primaryAlive:= true
	count := 0
	startCount := 0
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	checkErr(err)
	sConn, err := net.ListenUDP("udp", udpAddr)
	checkErr(err)

	primaryAlive = backupProcess(sConn, primaryAlive, &count)
	sConn.Close()

	fmt.Println("primary Process")
	startCount = count

	rAddr, _ := net.ResolveUDPAddr("udp4",Addr+port)
	mConn, err := net.DialUDP("udp4", nil, rAddr)
	if err != nil{
			fmt.Println("Error:connecting: ",err.Error())
	}


	
	for {
		msg := "Count:" + strconv.Itoa(count)
		_, err := mConn.Write([]byte(msg))
		fmt.Println("count: ", count)
		if err != nil {
			fmt.Println("Error:Broadcast", err.Error())
		}

		if (count == 5 + startCount) {
			break
		}
		count++
		time.Sleep(time.Second)
	}
	mConn.Close()
	fmt.Println("Finished counting")
	spawnProcess()

}
