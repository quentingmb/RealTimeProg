package main
 
import (
    "fmt"
    "net"
    "time"
    "strconv"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
 
func main() {
    ServerAddr,err := net.ResolveUDPAddr("udp","129.241.187.159:20011")
    CheckError(err)
 
    LocalAddr, err := net.ResolveUDPAddr("udp", "129.241.187.159:30000")
    CheckError(err)
 
    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)
 
    defer Conn.Close()
    i := 0
    for {
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := Conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second * 1)
        
        n, addr,err := Conn.ReadFromUDP(buf)
        
        fmt.Println("Received ", string(buf[0:n]), " from ",addr)
            
        }
}
