package multicast

import "strconv"

type DevLookingIntoBuild struct {
	DevName      string `json:"name"`
	PipelineName string `json:"pipeline_name"`
	JobName      string `json:"job_name"`
	ID           int `json:"ID"`
	Ignore       bool `json:"ignore"`
}

func (devReq DevLookingIntoBuild) Key() string {
	return strconv.Itoa(devReq.ID)
}