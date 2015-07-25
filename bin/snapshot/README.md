slackd
======

slackd is a very simple deamon that simply watches a text file (usually a log file) and pushes any lines which match a specific criteria to a [Slack](https://slack.com/) channel.  slackd allows you to specify a `line_includes` and/or a `line_excludes` regular expression to watch for, if both are provided they will both trigger a post to Slack.

slackd is a completely stand alone application and has no external dependances.  It does not need to be installed, it can simply be run inplace.


INSTALL
-------
```
$ git clone https://github.com/swill/slackd.git
$ cd slackd
$ tar zxfv ./bin/snapshot/slackd_[flavor].tar.gz -C . --strip-components 1
$ vim config.ini
    # create config according to usage section and save
$ nohup ./slackd -config=config.ini &
```


USAGE
-----
The application is self documented, so you can review the usage at any time.  

```
$ ./slackd -h
Usage of ./slackd:
  -channel="": The Slack channel to post to
  -config="": Path to ini config for using in go flags. May be relative to the current executable path.
  -file="": The file path to watch for changes
  -line_excludes="": Post if this regexp IS NOT in the line
  -line_includes="": Post if this regexp IS in the line
  -token="": Your Slack token
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
```
