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
	return fmt.Sprintf(`%s/%s %s Time Elapsed: %s | href=%s color=red
I will fix it | alternate=true bash=/usr/bin/say param1=test terminal=false`, build.PipelineName, build.JobName, build.Status, timeElapsed, targetUrl + build.URL)
}

func (p *Painter) Print() {
	fmt.Println(fmt.Sprintf("%d :red_circle: | color=red", len(p.MainItems)))
	fmt.Println("---")
	for _, value := range p.MainItems {
		fmt.Println(value)
	}
}