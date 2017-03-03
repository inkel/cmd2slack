# cmd2slack - Execute a command and send its output to Slack

Execute a command and send its `stdout` and `stderr` to a Slack channel as an [Incoming Webhook](https://api.slack.com/incoming-webhooks).

```
$ cmd2slack -hook https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX \
            cat /etc/hosts
```

## Installation

If you have a [Go](https://golang.org/) installation, just type `go get github.com/inkel/cmd2slack`.

## License

MIT. See [LICENSE](LICENSE).
