package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/concourse/fly/commands"
	"strings"
)

func main() {

	fly := &commands.Fly

	parser := flags.NewParser(fly, flags.HelpFlag | flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	iniParser := flags.NewIniParser(parser)
	iniParser.Parse(strings.NewReader(`
[Application Options]
; Concourse target name
Target = bosh

[builds]
Count = 50
`))
	fly.Builds.Execute(nil)
	fly.Pipelines.Execute(nil)

}