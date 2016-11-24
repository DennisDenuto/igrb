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
	command := fmt.Sprintf("bash=igrb param1=send param2=%s param3=a param4=b param5=c terminal=false", "dev-name", build.PipelineName, build.JobName, build.ID)
	return fmt.Sprintf(`%s/%s %s Time Elapsed: %s | href=%s color=red
I will fix it | alternate=true %s`, build.PipelineName, build.JobName, build.Status, timeElapsed, targetUrl + build.URL, command)
}

func (p *Painter) Print() {
	fmt.Println(fmt.Sprintf("%d :red_circle: | color=red", len(p.MainItems)))
	fmt.Println("---")
	for _, value := range p.MainItems {
		fmt.Println(value)
	}
}