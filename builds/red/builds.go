package red

import (
	"github.com/concourse/fly/rc"
	"fmt"
)

type Build struct {

}

type PersonInvestigating struct {

}


func ListBuilds(target string) {
	target, err := rc.LoadTarget(target)
	config, _, _, _, err := target.Team().PipelineConfig("bosh")

	for _, value := range config.Jobs {
		job, _, _ := target.Team().Job("bosh", value.Name)
		fmt.Println(job.Name)
		fmt.Println(job.FinishedBuild.Status)
	}
	fmt.Println(err)
}