package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type message struct {
	Text        string       `json:"text"`
	Username    string       `json:"username,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
}

type attachment struct {
	Fallback string   `json:"fallback,omitempty"`
	Color    string   `json:"color,omitempty"`
	Pretext  string   `json:"pretext,omitempty"`
	Text     string   `json:"text,omitempty"`
	MrkdwnIn []string `json:"mrkdwn_in,omitempty"`
	Fields   []field  `json:"fields,omitempty"`
	Ts       int      `json:"ts,omitempty"`
}

type field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func main() {
	msg := message{}

	hook := flag.String("hook", "", "Slack Incoming Webhook URL")
	flag.StringVar(&msg.Channel, "channel", "", "Channel where to post the output")
	flag.StringVar(&msg.IconEmoji, "emoji", "", "Emoji to use")
	flag.StringVar(&msg.Username, "username", "", "Username")
	flag.StringVar(&msg.IconURL, "icon", "", "URL of icon to use")
	flag.Parse()

	if *hook == "" {
		fmt.Fprintln(os.Stderr, "-hook is required")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "expected a command")
		os.Exit(2)
	}

	exe, args := args[0], args[1:]

	out, err := exec.Command(exe, args...).CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	msg.Text = "```\n" + string(out) + "\n```"

	body := new(bytes.Buffer)

	json.NewEncoder(body).Encode(msg)

	res, err := http.Post(*hook, "application/json", body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "posting to slack failed with", res.Status)
		// TODO this should be more expressive
		os.Exit(5)
	}
}
