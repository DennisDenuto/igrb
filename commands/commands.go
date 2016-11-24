package commands

import (
	multicast_reader "github.com/DennisDenuto/igrb/multicast/reader"
	"github.com/DennisDenuto/igrb/multicast/sender"
	"github.com/DennisDenuto/igrb/multicast"
	"github.com/DennisDenuto/igrb/data/diskstore"
)

const (
	srvAddr = "224.0.0.1:9999"
)

type ActionOpts struct {
	MulticastListen MulticastListenCommand `command:"listen"`
	MulticastSend   MulticastSendCommand `command:"send"`
	Status          StatusCommand `command:"status"`
	Ignore          IgnoreCommand `command:"ignore"`
}

type MulticastListenCommand struct{}
type StatusCommand struct{}
type MulticastSendCommand struct {
	Arg multicast.DevLookingIntoBuild `positional-args:"yes" required:"4"`
}
type IgnoreCommand struct{
	Arg multicast.DevLookingIntoBuild `positional-args:"yes" required:"4"`
}

func (MulticastListenCommand) Execute(args []string) error {
	multicast_reader.NewServeMulticastUDP(srvAddr).ServeMulticastUDP(nil)
	return nil
}

func (send MulticastSendCommand) Execute(args []string) error {
	return sender.NewMultiCastSender(srvAddr).SendMulticast(send.Arg)
}

func (ignore IgnoreCommand) Execute(args []string) error {
	ignore.Arg.Ignore = true
	return diskstore.NewDiskPersistor().Save(ignore.Arg.Key(), ignore.Arg)
}