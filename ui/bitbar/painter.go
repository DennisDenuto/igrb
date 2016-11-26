package bitbar

import (
	"fmt"
	"github.com/concourse/atc"
	"time"
)

type Painter struct {
	MainItems []string
}

func (p *Painter) AddMainMenuItems(item string) {
	p.MainItems = append(p.MainItems, item)
}

func JobToString(targetUrl string, build atc.Build) string {
	timeElapsed := time.Now().Sub(time.Unix(build.EndTime, 0))

	commandToInvestigate := fmt.Sprintf("bash=igrb param1=send param2=%s param3=%s param4=%s param5=%d terminal=true", "dev-name", build.PipelineName, build.JobName, build.ID)
	commandToIgnore := fmt.Sprintf("bash=igrb param1=ignore param2=%s param3=%s param4=%s param5=%d terminal=true", "_", build.PipelineName, build.JobName, build.ID)

	return fmt.Sprintf(`---
:exclamation: %s/%s %s | href=%s
--I got it! | %s
--Ignore | %s
Time red: %s`, build.PipelineName, build.JobName, build.Status, targetUrl + build.URL, commandToInvestigate, commandToIgnore, timeElapsed)
}

func (p *Painter) Print() {
	fmt.Println(fmt.Sprintf("%d :red_circle: | color=red", len(p.MainItems)))
	fmt.Println("---")
	for _, value := range p.MainItems {
		fmt.Println(value)
	}
}