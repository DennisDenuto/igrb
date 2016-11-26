package commands

import (
	"github.com/concourse/atc"
	"time"
)

type FailedBuildsSummary struct {
	URL          string `json:"url"`
	FailedBuilds map[string][]atc.Build `json:"failed_builds"`
	CreatedAt    time.Time `json:"created_at"`
}