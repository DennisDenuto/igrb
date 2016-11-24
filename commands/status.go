package commands

import (
	"fmt"
	"github.com/DennisDenuto/igrb/ui/bitbar"
	"github.com/concourse/fly/rc"
	"github.com/concourse/atc"
	"sync"
	"github.com/DennisDenuto/igrb/builds/red"
	"github.com/concourse/fly/commands"
	"strings"
	"github.com/jessevdk/go-flags"
)

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
