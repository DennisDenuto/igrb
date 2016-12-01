package multicast

import (
	"strconv"
	"time"
)

type DevLookingIntoBuild struct {
	DevName      string `json:"name"`
	PipelineName string `json:"pipeline_name"`
	JobName      string `json:"job_name"`
	ID           int `json:"ID"`
	Ignore       bool `json:"ignore"`
	CreatedAt    time.Time`json:"created_at"`
}

func (devReq DevLookingIntoBuild) Key() string {
	return strconv.Itoa(devReq.ID)
}