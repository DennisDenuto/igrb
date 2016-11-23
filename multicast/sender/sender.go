package sender

import (
	"net"
	"github.com/DennisDenuto/igrb/multicast"
	"encoding/json"
	logger "github.com/Sirupsen/logrus"
)

const (
	SrvAddr = "224.0.0.1:9999"
	MaxDatagramSize = 8192
)

type MulticastSender struct {
	SrvAddress string
}

func NewMultiCastSender(addr string) MulticastSender {
	return MulticastSender{
		SrvAddress: addr,
	}
}

func (sender MulticastSender) SendMulticast(devLookingIntoBuild multicast.DevLookingIntoBuild) error {
	addr, err := net.ResolveUDPAddr("udp", sender.SrvAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		logger.Error(err)
		return err
	}

	devJson, err := json.Marshal(devLookingIntoBuild)
	if err != nil {
		logger.Error(err)
		return err
	}

	c.Write(devJson)
	logger.Debugf("Sent %s", devJson)

	return nil
}