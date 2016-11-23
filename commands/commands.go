package commands

import "fmt"
import (
	multicast_reader "github.com/DennisDenuto/igrb/multicast/reader"
	"strings"
	"github.com/concourse/fly/rc"
	"sync"
	"github.com/concourse/atc"
	"github.com/DennisDenuto/igrb/ui/bitbar"
	"github.com/concourse/fly/commands"
	"github.com/jessevdk/go-flags"
	"github.com/DennisDenuto/igrb/builds/red"
	"github.com/DennisDenuto/igrb/multicast/sender"
	"github.com/DennisDenuto/igrb/multicast"
)

const (
	srvAddr = "224.0.0.1:9999"
)

type ActionOpts struct {
	MulticastListen MulticastListenCommand `command:"listen"`
	MulticastSend   MulticastSendCommand `command:"send"`
	Status          StatusCommand `command:"status"`
}

type MulticastListenCommand struct{}

type MulticastSendCommand struct {
	Arg multicast.DevLookingIntoBuild `positional-args:"yes" required:"4"`
}

type StatusCommand struct{}

func (MulticastListenCommand) Execute(args []string) error {
	var finish chan bool
	multicast_reader.NewServeMulticastUDP(srvAddr).ServeMulticastUDP(finish)

	return nil
}

func (send MulticastSendCommand) Execute(args []string) error {
	return sender.NewMultiCastSender(srvAddr).SendMulticast(send.Arg)
}

func (StatusCommand) Execute(args []string) error {
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
		return err
	}

	var wg sync.WaitGroup
	var failedBuilds map[string][]atc.Build = make(map[string][]atc.Build)

	var pipelines []atc.Pipeline
	pipelines, err = ListPipelines(fly.Pipelines.All, target)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, pipeline := range pipelines {
		wg.Add(1)
		go FetchFailedBuilds(pipeline, target, &wg, failedBuilds)
	}
	wg.Wait()

	painter := &bitbar.Painter{}
	for _, failedPipelineBuilds := range failedBuilds {
		for _, value := range failedPipelineBuilds {
			painter.AddMainMenuItems(bitbar.JobToString(target.URL(), value))
		}
	}

	painter.Print()
	return nil
}

func ListPipelines(all bool, target rc.Target) ([]atc.Pipeline, error) {
	if all {
		return target.Client().ListPipelines()
	} else {
		return target.Team().ListPipelines()
	}
}

func FetchFailedBuilds(pipeline atc.Pipeline, target rc.Target, wg *sync.WaitGroup, failedBuilds map[string][]atc.Build) {
	defer wg.Done()

	failedBuildsForPipeline, err := red.FailedBuildFetcher{Target: target}.Fetch(pipeline.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	failedBuilds[pipeline.Name] = failedBuildsForPipeline
}
