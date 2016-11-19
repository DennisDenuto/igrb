package main

import (
	"github.com/jessevdk/go-flags"
	"strings"
	"fmt"
	"github.com/DennisDenuto/igrb/multicast/sender"
	"github.com/concourse/fly/commands"
	"github.com/concourse/fly/rc"
	multicast_reader "github.com/DennisDenuto/igrb/multicast/reader"
	"github.com/DennisDenuto/igrb/builds/red"
)

const (
	srvAddr         = "224.0.0.1:9999"
)


func main() {
	go sender.Ping(srvAddr)
	multicast_reader .ServeMulticastUDP(srvAddr, multicast_reader.MsgHandler)
}

func main1() {

	fly := &commands.Fly

	parser := flags.NewParser(fly, flags.HelpFlag | flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	iniParser := flags.NewIniParser(parser)
	iniParser.Parse(strings.NewReader(`
[Application Options]
; Concourse target name
Target = bosh

[builds]
Count = 50
`))

	red.ListBuilds(fly.Target)
}