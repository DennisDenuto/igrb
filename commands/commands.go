package commands

import (
	multicast_reader "github.com/DennisDenuto/igrb/multicast/reader"
	"github.com/DennisDenuto/igrb/multicast/sender"
	"github.com/DennisDenuto/igrb/multicast"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"strings"
	"github.com/concourse/fly/rc"
	"github.com/concourse/atc"
	"time"
	"github.com/concourse/fly/commands"
	"github.com/jessevdk/go-flags"
	"fmt"
	"sync"
	"github.com/pkg/errors"
	logger "github.com/Sirupsen/logrus"
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
	Arg multicast.DevLookingIntoBuild `positional-args:"yes"`
}
type IgnoreCommand struct {
	Arg multicast.DevLookingIntoBuild `positional-args:"yes"`
}

func updatedFailedBuildSummary(concourseTarget string) error {
	fly := &commands.Fly

	parser := flags.NewParser(fly, flags.HelpFlag | flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	iniParser := flags.NewIniParser(parser)
	iniParser.Parse(strings.NewReader(fmt.Sprintf(`
[Application Options]
; Concourse target name
Target = %s

[builds]
Count = 50
`, concourseTarget)))

	target, err := rc.LoadTarget(fly.Target)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = fetchFailedBuildsRemotely(fly, target)
	return err
}

func (MulticastListenCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("Missing concourse target argument")
	}
	go func() {
		for {
			err := updatedFailedBuildSummary(args[0])
			if err != nil {
				logger.Error(errors.Wrap(err, "Unable to update failed build summary"))
			}
			time.Sleep(10 * time.Second)
		}
	}()

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

func fetchFailedBuildsRemotely(fly *commands.FlyCommand, target rc.Target) (map[string][]atc.Build, error) {
	var failedBuilds map[string][]atc.Build = make(map[string][]atc.Build)

	var wg sync.WaitGroup
	var pipelines []atc.Pipeline
	pipelines, err := ListPipelines(fly.Pipelines.All, target)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, pipeline := range pipelines {
		wg.Add(1)
		go FetchFailedBuilds(pipeline, target, &wg, failedBuilds)
	}
	wg.Wait()

	err = SaveFailedBuildsSummary(target.URL(), failedBuilds)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "Unable to save summary of failed builds")
	}

	return failedBuilds, nil

}
