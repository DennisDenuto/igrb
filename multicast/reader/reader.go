package reader

import (
	"net"
	"github.com/DennisDenuto/igrb/multicast"
	"log"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"encoding/json"
	"github.com/DennisDenuto/igrb/multicast/types"
)

type MulticastReceiver struct {
	SrvAddress string
	handler    func(src *net.UDPAddr, n int, b []byte)
}

func NewServeMulticastUDP(srvAddress string) MulticastReceiver {
	return MulticastReceiver{
		SrvAddress: srvAddress,
		handler: MsgHandler,
	}
}

func MsgHandler(src *net.UDPAddr, n int, b []byte) {
	devReq := types.DevLookingIntoBuild{}
	json.Unmarshal(b[:n], &devReq)

	diskstore.NewDiskPersistor().Save("pipeline_job_build-id", devReq)
}

func (receiver MulticastReceiver) ServeMulticastUDP(finish <-chan bool) {
	addr, err := net.ResolveUDPAddr("udp", receiver.SrvAddress)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(multicast.MaxDatagramSize)
	for {
		select {
		case <-finish:
			return
		default:
			b := make([]byte, multicast.MaxDatagramSize)
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
				log.Fatal("ReadFromUDP failed:", err)
			}
			receiver.handler(src, n, b)
		}
	}
}

