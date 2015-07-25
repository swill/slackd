package main

import (
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/nlopes/slack"
	"github.com/vharitonsky/iniflags"
	"os"
	"regexp"
)

var (
	token    = flag.String("token", "", "Your Slack token")
	channel  = flag.String("channel", "", "The Slack channel to post to (without the leading '#')")
	file     = flag.String("file", "", "The file path to watch for changes")
	includes = flag.String("line_includes", "", "Post line if this regexp DOES match")
	excludes = flag.String("line_excludes", "", "Post line if this regexp DOES NOT match")
)

func main() {
	iniflags.Parse()

	api := slack.New(*token)
	channels, err := api.GetChannels(true)
	if err != nil {
		fmt.Println("\nERROR: Could not get the Slack channels\n")
		fmt.Println(err)
		os.Exit(2)
	}
	var channel_id string
	for _, c := range channels {
		if c.Name == *channel {
			channel_id = c.Id
		}
	}
	if channel_id == "" {
		fmt.Println("\nERROR: Could not find the Slack channel specified.  Be sure NOT to include the '#' at the beginning.\n")
		os.Exit(2)
	}

	var include, exclude *regexp.Regexp
	if *includes != "" {
		include, err = regexp.Compile(*includes)
		if err != nil {
			fmt.Println("\nERROR: Failed to compile `line_includes` regex\n")
			fmt.Println(err)
			api.PostMessage(channel_id, "==> slackd failed to compile `line_includes` regex", slack.NewPostMessageParameters())
			api.PostMessage(channel_id, err.Error(), slack.NewPostMessageParameters())
			os.Exit(2)
		}
	}
	if *excludes != "" {
		exclude, err = regexp.Compile(*excludes)
		if err != nil {
			fmt.Println("\nERROR: Failed to compile `line_excludes` regex\n")
			fmt.Println(err)
			api.PostMessage(channel_id, "==> slackd failed to compile `line_excludes` regex", slack.NewPostMessageParameters())
			api.PostMessage(channel_id, err.Error(), slack.NewPostMessageParameters())
			os.Exit(2)
		}
	}

	log, err := tail.TailFile(*file, tail.Config{Follow: true})
	if err != nil {
		fmt.Println("\nERROR: Could not tail the specified log\n")
		fmt.Println(err)
		api.PostMessage(channel_id, "==> slackd could not tail the specified log", slack.NewPostMessageParameters())
		api.PostMessage(channel_id, err.Error(), slack.NewPostMessageParameters())
		os.Exit(2)
	}
	for line := range log.Lines {
		if (include != nil && include.MatchString(line.Text)) || (exclude != nil && !exclude.MatchString(line.Text)) {
			api.PostMessage(channel_id, fmt.Sprintf("`%s`", line.Text), slack.NewPostMessageParameters())
		}
	}
}
