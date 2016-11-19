package red

import (
	"github.com/concourse/fly/rc"
	"github.com/concourse/atc"
	"github.com/pkg/errors"
)

type FailedBuildFetcher struct {
	Target rc.Target
}

type PersonInvestigating struct {
}

func (build FailedBuildFetcher) Fetch(pipelineName string) ([]atc.Build, error) {
	team := build.Target.Team()
	config, _, _, _, err := team.PipelineConfig(pipelineName)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to get pipeline config")
	}
	var builds []atc.Build

	for _, value := range config.Jobs {
		job, _, err := team.Job(pipelineName, value.Name)

		if err != nil {
			return nil, errors.Wrap(err, "Unable to get pipeline config")
		}
		switch job.FinishedBuild.Status {
		case string(atc.StatusFailed): builds = append(builds, *job.FinishedBuild)
		case string(atc.StatusErrored): builds = append(builds, *job.FinishedBuild)
		}
	}

	return builds, nil
}