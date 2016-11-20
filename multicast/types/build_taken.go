package types

type DevLookingIntoBuild struct {
	DevName      string `json:"name"`
	PipelineName string `json:"build_name"`
	JobName      string `json:"job_name"`
	BuildId      string `json:"build_id"`
}

