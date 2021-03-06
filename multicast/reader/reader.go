package reader

import (
	"net"
	"github.com/DennisDenuto/igrb/multicast"
	"log"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"encoding/json"
	"github.com/pkg/errors"
	logger "github.com/Sirupsen/logrus"
	"github.com/DennisDenuto/igrb/multicast/sender"
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
	devReq := multicast.DevLookingIntoBuild{}
	err := json.Unmarshal(b[:n], &devReq)
	if err != nil {
		logger.Error(errors.Wrap(err, "Unable to parse multicast request"))
	}

	if devReq.PipelineName == "" || devReq.JobName == "" {
		logger.Info("Skipping devRequest due to missing fields")
		return
	}

	existingDevLookingIntoBuild := multicast.DevLookingIntoBuild{}
	err = diskstore.NewDiskPersistor().ReadAndUnmarshal(devReq.Key(), &existingDevLookingIntoBuild)

	if err == nil && existingDevLookingIntoBuild.CreatedAt.After(devReq.CreatedAt) {
		logger.Debugf("Skipping %s due to existing record being newer", devReq)
		return
	}

	logger.Debugf("Saving %s", devReq)
	diskstore.NewDiskPersistor().Save(devReq.Key(), devReq)
}

func (receiver MulticastReceiver) ServeMulticastUDP(finish <-chan bool) {
	addr, err := net.ResolveUDPAddr("udp", receiver.SrvAddress)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	logger.Debugf("Listening on %s", addr)
	l.SetReadBuffer(sender.MaxDatagramSize)
	for {
		select {
		case <-finish:
			return
		default:
			b := make([]byte, sender.MaxDatagramSize)
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
				log.Fatal("ReadFromUDP failed:", err)
			}
			receiver.handler(src, n, b)
		}
	}
}

