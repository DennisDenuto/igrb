package main

// <bitbar.title>I got the red build</bitbar.title>
// <bitbar.version>v1.0</bitbar.version>
// <bitbar.author>DennisDenuto</bitbar.author>
// <bitbar.author.github>DennisDenuto</bitbar.author.github>
// <bitbar.desc>Allows a developer to take on the responsibility for changing a red build -> green build</bitbar.desc>
// <bitbar.abouturl>https://github.com/DennisDenuto/igrb</bitbar.abouturl>
//
// Text above --- will be cycled through in the menu bar,
// whereas text underneath will be visible only when you
// open the menu.
//

import (
	"github.com/jessevdk/go-flags"
	"strings"
	"github.com/DennisDenuto/igrb/multicast/sender"
	"github.com/concourse/fly/commands"
	multicast_reader "github.com/DennisDenuto/igrb/multicast/reader"
	"github.com/DennisDenuto/igrb/builds/red"
	"fmt"
	"github.com/concourse/fly/rc"
	"github.com/DennisDenuto/igrb/ui/bitbar"
	"sync"
	"github.com/concourse/atc"
)

const (
	srvAddr = "224.0.0.1:9999"
)

func main1() {
	go sender.Ping(srvAddr)
	multicast_reader.ServeMulticastUDP(srvAddr, multicast_reader.MsgHandler)
}

func main() {

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
	target, err := rc.LoadTarget(fly.Target)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	var pipelines []atc.Pipeline
	if fly.Pipelines.All {
		pipelines, err = target.Client().ListPipelines()
	} else {
		pipelines, err = target.Team().ListPipelines()
	}

	var failedBuilds map[string][]atc.Build = make(map[string][]atc.Build)

	for _, p := range pipelines {
		wg.Add(1)
		go func(pipeline atc.Pipeline) {
			defer wg.Done()

			failedBuildsForPipeline, err := red.FailedBuildFetcher{Target: target}.Fetch(pipeline.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
			failedBuilds[pipeline.Name] = failedBuildsForPipeline
		}(p)
	}
	wg.Wait()

	painter := &bitbar.Painter{}
	for _, failedPipelineBuilda := range failedBuilds {
		for _, value := range failedPipelineBuilda {
			painter.AddMainMenuItems(bitbar.JobToString(target.URL(), value))
		}
	}

	painter.Print()
}