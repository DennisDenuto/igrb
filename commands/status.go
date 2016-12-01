package commands

import (
	"fmt"
	"github.com/DennisDenuto/igrb/ui/bitbar"
	"github.com/concourse/fly/rc"
	"github.com/concourse/atc"
	"sync"
	"github.com/DennisDenuto/igrb/builds/red"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"github.com/DennisDenuto/igrb/multicast"
	"strconv"
	"time"
	"sort"
)

const SUMMARY_KEY = "summary"

func (StatusCommand) Execute(args []string) error {
	var failedBuilds map[string][]atc.Build

	summary := &FailedBuildsSummary{}
	diskstore.NewDiskPersistor().ReadAndUnmarshal(SUMMARY_KEY, summary)
	failedBuilds = summary.FailedBuilds

	painter := &bitbar.Painter{}

	AddMenuItemsToPainter(summary.URL, failedBuilds, painter)
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
	var failedBuildsNotIgnored []atc.Build

	for _, failedPipelineBuilds := range failedBuilds {
		for _, value := range failedPipelineBuilds {
			if buildIgnored(value) {
				continue
			}
			failedBuildsNotIgnored = append(failedBuildsNotIgnored, value)
		}
	}

	sort.Sort(atc.Builds(failedBuildsNotIgnored))
	for _, value := range failedBuildsNotIgnored {
		painter.AddMainMenuItems(bitbar.JobToString(url, value))
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
