package main

// <bitbar.title>I got the red build</bitbar.title>
// <bitbar.version>v1.0</bitbar.version>
// <bitbar.author>DennisDenuto</bitbar.author>
// <bitbar.author.github>DennisDenuto</bitbar.author.github>
// <bitbar.desc>Allows a developer to take on the responsibility for changing a red build -> green build</bitbar.desc>
// <bitbar.abouturl>https://github.com/DennisDenuto/igrb</bitbar.abouturl>
//
// Text above --- will be cycled through in the menu bar,
// whereas text underneath will be visible only when you
// open the menu.
//

import (
	"os"

	"github.com/DennisDenuto/igrb/commands"
	logger "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
	var f *os.File
	var err error
	if f, err = os.OpenFile("/tmp/igrb.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0750); err != nil {
		logger.Error("create log failed", err)
		return
	}
	logger.SetOutput(f)
	defer func(f *os.File) {
		f.Sync()
	}(f)

	logger.SetLevel(logger.DebugLevel)
}

func main() {
	actionOpts := commands.ActionOpts{}
	flags.Parse(&actionOpts)
}
