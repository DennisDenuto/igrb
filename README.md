# igrb
I got the red build - OSX concourse bitbar plugin for a team to communicate who is working on a red build

igrb auto discovers peers on the network using multicast.

![screenshot](igrb.png)

# Usage
* Install go 1.7
* Install bitbar (`brew install bitbar`)
* Install igrb (`go install github.com/DennisDenuto/igrb`)
* Add a script like the following to `~/.bitbar/` and make it executable:

  ```sh
#!/bin/bash
CONCOURSE_TARGET=<your concourse target>

if ! pgrep -qf "listen $CONCOURSE_TARGET"; then
  nohup /usr/local/bin/igrb listen $CONCOURSE_TARGET &
fi
/usr/local/bin/igrb status $CONCOURSE_TARGET
  ```

  
** Since the bitbar protocol is text-based, the bitbar-igrb plugin can be tested in the terminal. Just execute the script in `~/.bitbar/` in a terminal window.
