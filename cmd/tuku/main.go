package main

import (
	"github.com/alecthomas/kong"
	"github.com/sinkingpoint/kiora/cmd/tuku/commands"
	"github.com/sinkingpoint/kiora/cmd/tuku/commands/alerts"
	"github.com/sinkingpoint/kiora/cmd/tuku/kiora"
	"github.com/sinkingpoint/kiora/internal/encoding"
)

var CLI struct {
	Formatter string           `help:"the format to output the data in" default:"json"`
	KioraURL  string           `help:"the URL of the Kiora instance to connect to" default:"http://localhost:4278"`
	Alerts    alerts.AlertsCmd `cmd:"" help:"Manage alerts."`
}

func main() {
	ctx := kong.Parse(&CLI, kong.Name("tuku"), kong.Description("A CLI for interacting with Kiora"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{
		Compact: true,
	}))

	runContext := &commands.Context{
		Formatter: encoding.LookupEncoding(CLI.Formatter),
		Kiora:     kiora.NewKioraInstance(CLI.KioraURL, "v1"),
	}

	if err := ctx.Run(runContext); err != nil {
		ctx.FatalIfErrorf(err)
	}
}
