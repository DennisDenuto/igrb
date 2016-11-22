package multicast

import (
	"fmt"
	"strings"
)

type DevLookingIntoBuild struct {
	DevName      string `json:"name"`
	PipelineName string `json:"build_name"`
	JobName      string `json:"job_name"`
	BuildId      string `json:"build_id"`
}

func (devReq DevLookingIntoBuild) Key() string {
	return fmt.Sprintf("%s_%s_%s", strings.Replace(devReq.PipelineName, "/", "_", -1), strings.Replace(devReq.JobName, "/", "_", -1), devReq.BuildId)
}