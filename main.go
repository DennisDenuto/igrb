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
	"github.com/jessevdk/go-flags"
	"github.com/DennisDenuto/igrb/commands"
	"os"
	logger "github.com/Sirupsen/logrus"
)


func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
	logger.SetOutput(os.Stderr)
	logger.SetLevel(logger.DebugLevel)
}


func main() {
	actionOpts := commands.ActionOpts{}
	flags.Parse(&actionOpts)
}
