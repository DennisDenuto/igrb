# igrb
I got the red build - OSX concourse bitbar plugin for a team to communicate who is working on a red build

igrb auto discovers peers on the network using multicast.

![screenshot](igrb.png)

# Usage
* Install go 1.7
* Install bitbar (`brew install Caskroom/cask/bitbar`)
* Install igrb (`go get github.com/DennisDenuto/igrb`)
* Add a script like the following to `~/.bitbar/igrb.10s.sh` and make it executable:

```sh
#!/bin/bash
CONCOURSE_TARGET=<your concourse target>

if ! pgrep -qf "listen $CONCOURSE_TARGET"; then
  nohup /usr/local/bin/igrb listen $CONCOURSE_TARGET &
fi
/usr/local/bin/igrb status $CONCOURSE_TARGET
```

  
# Configure the refresh time

The refresh time is in the filename of the plugin, following this format:

{name}.{time}.{ext}
name - The name of the file
time - The refresh rate (see below)
ext - The file extension
For example:

igrb.10s.sh would refresh every 10 seconds.

** Since the bitbar protocol is text-based, the bitbar-igrb plugin can be tested in the terminal. Just execute the script in `~/.bitbar/` in a terminal window.
