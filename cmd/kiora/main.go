package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/internal/server"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

var CLI struct {
	HTTPListenAddress string `name:"web.listen-url" help:"the address to listen on" default:"localhost:4278"`
	ConfigFile        string `name:"config.file" short:"c" help:"the config file to load config from" default:"./kiora.dot"`

	NodeName             string   `name:"cluster.node-name" help:"the name to join the cluster with"`
	ClusterListenAddress string   `name:"cluster.listen-url" help:"the address to run cluster activities on" default:"localhost:4279"`
	BootstrapPeers       []string `name:"cluster.bootstrap-peers" help:"the peers to bootstrap with"`
}

func main() {
	kong.Parse(&CLI, kong.Name("kiora"), kong.Description("An experimental Alertmanager"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{
		Compact: true,
	}))

	config, err := config.LoadConfigFile(CLI.ConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	serverConfig := server.NewServerConfig()
	serverConfig.HTTPListenAddress = CLI.HTTPListenAddress
	serverConfig.ClusterListenAddress = CLI.ClusterListenAddress
	serverConfig.BootstrapPeers = CLI.BootstrapPeers
	serverConfig.NotifierConfig = config

	tracingConfig := tracing.DefaultTracingConfiguration()
	tracingConfig.ExporterType = "jaeger" // TODO(cdouch): Make this a CLI arg
	tp, err := tracing.InitTracing(tracingConfig)
	if err != nil {
		log.Warn().Err(err).Msg("failed to start tracing")
	}

	if tp != nil {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Warn().Err(err).Msg("failed to shutdown tracing. Spans may have been lost")
			}
		}()
	}

	server, err := server.NewKioraServer(serverConfig, kioradb.NewInMemoryDB())
	if err != nil {
		log.Err(err).Msg("failed to create server")
		return
	}

	// Setup a SIGINT handler, so that we can shutdown gracefully.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Info().Msg("Received signal, shutting down")
			server.Shutdown()
			break
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("failed to listen and serve")
	}

	log.Info().Msg("Kiora Shut Down")
}
