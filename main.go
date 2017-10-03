package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
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

func (a *attachment) addField(title, value string) {
	a.Fields = append(a.Fields, field{Title: title, Value: value, Short: true})
}

func main() {
	msg := message{}

	var (
		hook    = flag.String("hook", "", "Slack Incoming Webhook URL")
		timing  = flag.Bool("timing", false, "Include command execution timing")
		verbose = flag.Bool("verbose", false, "Show command execution on screen")
	)

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

	a := attachment{
		Color:    "good",
		Fallback: fmt.Sprintf("%s results", exe),
		Pretext:  fmt.Sprintf("`%s`", pretty(exe, args)),
		MrkdwnIn: []string{"pretext", "text"},
	}

	start := time.Now()

	out := bytes.NewBuffer(nil)

	cmd := exec.Command(exe, args...)
	cmd.Stdout = out
	cmd.Stderr = out

	if *verbose {
		cmd.Stdout = io.MultiWriter(out, os.Stdout)
		cmd.Stderr = io.MultiWriter(out, os.Stderr)
	}

	err := cmd.Run()
	if err != nil {
		a.Color = "danger"
		a.addField("Error", err.Error())
		fmt.Fprintln(os.Stderr, err)
	}
	if out.Len() > 0 {
		a.Text = fmt.Sprintf("```\n%s```", out)
	}

	if *timing {
		a.addField("Timing", time.Since(start).String())
	}

	msg.Attachments = []attachment{a}

	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	res, err := http.Post(*hook, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "posting to slack failed with", res.Status)
		io.Copy(os.Stderr, res.Body)
		os.Exit(5)
	}
}

func pretty(cmd string, args []string) string {
	buf := bytes.NewBufferString(cmd)

	for _, arg := range args {
		fmts := " %s"
		for _, b := range arg {
			if isSpace(b) {
				fmts = " %q"
				break
			}
		}

		fmt.Fprintf(buf, fmts, arg)
	}
	return buf.String()
}

// Borrowed from unicode.IsSpace
func isSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}
