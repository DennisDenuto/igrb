package multicast

import "github.com/DennisDenuto/igrb/builds/red"

const (
	SrvAddr         = "224.0.0.1:9999"
	MaxDatagramSize = 8192
)


type Sender interface {
	Send(red.FailedBuildFetcher, red.PersonInvestigating) error
}
