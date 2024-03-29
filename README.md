slackd
======

slackd is a very simple daemon that simply watches a text file (usually a log file) and pushes any lines which match a specific criteria to a [Slack](https://slack.com/) channel.  slackd allows you to specify a `line_includes` and/or a `line_excludes` regular expression to watch for, if both are provided they will both trigger a post to Slack.

slackd is a completely stand alone application and has no external dependencies.  It does not need to be installed, it can simply be run in-place.


INSTALL
-------
In this example I am using the `slackd_linux_amd64` binary, but you can change it to whichever flavor you need to use.

**Using `rc.local`**
```
$ wget -O slackd https://github.com/swill/slackd/raw/master/bin/slackd_linux_amd64
$ vim config.ini
    # create config according to usage section and save
$ sudo vim /etc/rc.local
	# add the following before the `exit 0`
	# cd /path/to/slackd && nohup ./slackd -config=config.ini &
$ nohup ./slackd -config=config.ini &
```

**Using `systemd`**
```
$ cp slackd.service.example /etc/systemd/system/slackd.service
```

Edit the environment variables and location of the binary to fit your needs.

Start Service
```
sudo systemctl start slackd
```

Stop Service
```
sudo systemctl stop slackd
```

View Logs
```
journalctl -u slackd
```

USAGE
-----
The application is self documented, so you can review the usage at any time.

**Note:** All uppercase flags can be set via environment variables in addition to the config file and the CLI.

```
$ ./slackd -h
Usage of slackd:
  -CHANNEL string
      The ID of the Slack channel to post to. EG: C0XXXXXXXXX
  -config string
      Path to ini config for using in go flags. May be relative to the current executable path.
  -configUpdateInterval duration
      Update interval for re-reading config file set via -config flag. Zero disables config file re-reading.
  -dumpflags
      Dumps values for all flags defined in the app into stdout in ini-compatible syntax and terminates the app.
  -FILE string
      The file path to watch for changes.
  -LINE_EXCLUDES string
      Post line if this regexp DOES NOT match.
  -LINE_INCLUDES string
      Post line if this regexp DOES match.
  -REOPEN
      Reopen the file if it disappears. Useful with logrotation.
  -TOKEN string
      Your Slack token.
```


EXAMPLE CONFIG
--------------

With this config I want to watch the `/var/log/application.log` file.  I want to post lines from the log file to the Slack channel `errors` for either of the following conditions:

- The line includes the case insensitive text `error`.
- The line does not start with a date in the specified format (useful for stack traces).

```
token = ####-##########-##########-##########-######
channel = errors
file = /var/log/application.log
line_includes = (?i)error
line_excludes = ^[0-9]{4}/[0-9]{2}/[0-9]{2}
reopen = true
```
