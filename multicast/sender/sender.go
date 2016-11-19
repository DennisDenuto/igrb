package sender
import (
	"log"
	"net"
	"time"
)


func Ping(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	for {
		c.Write([]byte("(((((((hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}
