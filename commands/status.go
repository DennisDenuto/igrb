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
	"github.com/DennisDenuto/igrb/data/diskstore"
	"github.com/DennisDenuto/igrb/multicast"
	"strconv"
	"time"
	"github.com/pkg/errors"
)

const SUMMARY_KEY = "summary"

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

	err = SaveFailedBuildsSummary(target.URL(), failedBuilds)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "Unable to save summary of failed builds")
	}

	painter := &bitbar.Painter{}

	AddMenuItemsToPainter(target.URL(), failedBuilds, painter)
	painter.Print()

	return nil
}

func SaveFailedBuildsSummary(url string, failedBuilds map[string][]atc.Build) error {
	summary := FailedBuildsSummary{
		URL: url,
		FailedBuilds: failedBuilds,
		CreatedAt: time.Now(),
	}

	return diskstore.NewDiskPersistor().Save(SUMMARY_KEY, summary)
}

func AddMenuItemsToPainter(url string, failedBuilds map[string][]atc.Build, painter *bitbar.Painter) {
	for _, failedPipelineBuilds := range failedBuilds {
		for _, value := range failedPipelineBuilds {
			if buildIgnored(value) {
				continue
			}
			painter.AddMainMenuItems(bitbar.JobToString(url, value))
		}
	}
}

func buildIgnored(build atc.Build) bool {
	buildReq := &multicast.DevLookingIntoBuild{}
	diskstore.NewDiskPersistor().ReadAndUnmarshal(strconv.Itoa(build.ID), buildReq)

	return buildReq.Ignore
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
