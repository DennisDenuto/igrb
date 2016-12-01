package bitbar

import (
	"fmt"
	"github.com/concourse/atc"
	"time"
	"github.com/DennisDenuto/igrb/data/diskstore"
	"github.com/DennisDenuto/igrb/multicast"
	"strconv"
	"github.com/git-duet/git-duet"
)

type Painter struct {
	MainItems []string
}

func (p *Painter) AddMainMenuItems(item string) {
	p.MainItems = append(p.MainItems, item)
}

func JobToString(targetUrl string, build atc.Build) string {
	timeElapsed := time.Now().Sub(time.Unix(build.EndTime, 0))

	commandToInvestigate := fmt.Sprintf("bash=/usr/local/bin/igrb param1=send param2=%s param3=%s param4=%s param5=%d terminal=false refresh=true", GetGitUser(), build.PipelineName, build.JobName, build.ID)
	commandToIgnore := fmt.Sprintf("bash=/usr/local/bin/igrb param1=ignore param2=%s param3=%s param4=%s param5=%d terminal=false refresh=true", "_", build.PipelineName, build.JobName, build.ID)

	var icon string = ":exclamation:"
	buildTakenByDev := &multicast.DevLookingIntoBuild{}
	diskstore.NewDiskPersistor().ReadAndUnmarshal(strconv.Itoa(build.ID), buildTakenByDev)
	if buildTakenByDev.DevName != "" {
		icon = ":grey_question:"
	}

	buildSummaryText := fmt.Sprintf(`---
%s %s/%s %s | href=%s
--I got it! | %s
--Ignore | %s
Time red: %s`, icon, build.PipelineName, build.JobName, build.Status, targetUrl + build.URL, commandToInvestigate, commandToIgnore, timeElapsed)

	var buildFooter string

	if buildTakenByDev.DevName != "" {
		buildFooter = fmt.Sprintf("%s is looking into it!", buildTakenByDev.DevName)
	}

	return buildSummaryText + "\n" + buildFooter
}

func (p *Painter) Print() {
	fmt.Println(fmt.Sprintf("%d :red_circle: | color=red", len(p.MainItems)))
	fmt.Println("---")
	for _, value := range p.MainItems {
		fmt.Println(value)
	}
}

func GetGitUser() string {
	configuration, err := duet.NewConfiguration()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	gitConfig := &duet.GitConfig{Namespace: "user", SetUserConfig: configuration.SetGitUserConfig}
	gitDuetConfig := &duet.GitConfig{Namespace: configuration.Namespace, SetUserConfig: configuration.SetGitUserConfig}
	pair, err := gitDuetConfig.GetAuthor()
	if err != nil {
		return pair.Name
	}
	name, err := gitConfig.GetKey("name")
	return name
}