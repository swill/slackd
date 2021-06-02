package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/hpcloud/tail"
	"github.com/nlopes/slack"
	"github.com/vharitonsky/iniflags"
)

var (
	token    string
	channel  string
	file     string
	includes string
	excludes string
	reopen   bool
	TOKEN    = flag.String("TOKEN", "", "Your Slack token.")
	CHANNEL  = flag.String("CHANNEL", "", "The ID of the Slack channel to post to. EG: C0XXXXXXXXX")
	FILE     = flag.String("FILE", "", "The file path to watch for changes.")
	INCLUDES = flag.String("LINE_INCLUDES", "", "Post line if this regexp DOES match.")
	EXCLUDES = flag.String("LINE_EXCLUDES", "", "Post line if this regexp DOES NOT match.")
	REOPEN   = flag.Bool("REOPEN", false, "Reopen the file if it disappears. Useful with logrotation.")
)

func main() {
	iniflags.Parse()

	token = *TOKEN
	if token == "" && os.Getenv("TOKEN") != "" {
		token = os.Getenv("TOKEN")
	}
	channel = *CHANNEL
	if channel == "" && os.Getenv("CHANNEL") != "" {
		channel = os.Getenv("CHANNEL")
	}
	file = *FILE
	if file == "" && os.Getenv("FILE") != "" {
		file = os.Getenv("FILE")
	}
	includes = *INCLUDES
	if includes == "" && os.Getenv("LINE_INCLUDES") != "" {
		includes = os.Getenv("LINE_INCLUDES")
	}
	excludes = *EXCLUDES
	if excludes == "" && os.Getenv("LINE_EXCLUDES") != "" {
		excludes = os.Getenv("LINE_EXCLUDES")
	}
	reopen = *REOPEN
	if reopen == false && os.Getenv("REOPEN") != "" {
		var err error
		reopen, err = strconv.ParseBool(os.Getenv("REOPEN"))
		if err != nil {
			fmt.Println("ERROR: Parsing the REOPEN boolean.")
			fmt.Println(err)
			os.Exit(2)
		}
	}

	slackAPI := slack.New(token)
	var include, exclude *regexp.Regexp
	var err error
	if includes != "" {
		include, err = regexp.Compile(includes)
		if err != nil {
			fmt.Println("ERROR: Failed to compile `LINE_INCLUDES` regex.")
			fmt.Println(err)
			opt := slack.MsgOptionText(fmt.Sprintf("==> slackd failed to compile `LINE_INCLUDES` regex."), false)
			slackAPI.PostMessage(channel, opt)
			os.Exit(2)
		}
	}
	if excludes != "" {
		exclude, err = regexp.Compile(excludes)
		if err != nil {
			fmt.Println("ERROR: Failed to compile `LINE_EXCLUDES` regex.")
			fmt.Println(err)
			opt := slack.MsgOptionText(fmt.Sprintf("==> slackd failed to compile `LINE_EXCLUDES` regex."), false)
			slackAPI.PostMessage(channel, opt)
			os.Exit(2)
		}
	}

	log, err := tail.TailFile(file, tail.Config{Follow: true, ReOpen: reopen, Poll: true})
	if err != nil {
		fmt.Println("ERROR: Could not tail the specified log.")
		fmt.Println(err)
		opt := slack.MsgOptionText(fmt.Sprintf("==> slackd could not tail the specified log."), false)
		slackAPI.PostMessage(channel, opt)
		os.Exit(2)
	}
	for line := range log.Lines {
		if (include != nil && include.MatchString(line.Text)) || (exclude != nil && !exclude.MatchString(line.Text)) {
			opt := slack.MsgOptionText(fmt.Sprintf("```%s```", line.Text), false)
			slackAPI.PostMessage(channel, opt)
		}
	}
}
