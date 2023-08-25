package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/internal/server"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

var CLI struct {
	tracing.TracingConfiguration ` prefix:"tracing."`
	HTTPListenAddress            string `name:"web.listen-url" help:"the address to listen on" default:"localhost:4278"`
	ConfigFile                   string `name:"config.file" short:"c" help:"the config file to load config from" default:"./kiora.dot"`

	NodeName             string   `name:"cluster.node-name" help:"the name to join the cluster with"`
	ClusterListenAddress string   `name:"cluster.listen-url" help:"the address to run cluster activities on" default:"localhost:4279"`
	ClusterShardLabels   []string `name:"cluster.shard-labels" help:"the labels that determine which node in a cluster will send a given alert"`
	BootstrapPeers       []string `name:"cluster.bootstrap-peers" help:"the peers to bootstrap with"`

	StorageBackend string `name:"storage.backend" help:"the storage backend to use" default:"boltdb"`
	StoragePath    string `name:"storage.path" help:"the path to store data in" default:"./kiora.db"`
}

func main() {
	CLI.TracingConfiguration = tracing.DefaultTracingConfiguration()
	kong.Parse(&CLI, kong.Name("kiora"), kong.Description("An experimental Alertmanager"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{
		Compact: true,
	}))

	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
	log.Logger = logger

	config.RegisterNodes()
	config, err := config.LoadConfigFile(CLI.ConfigFile, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	serverConfig := server.NewServerConfig()
	serverConfig.HTTPListenAddress = CLI.HTTPListenAddress
	serverConfig.ClusterListenAddress = CLI.ClusterListenAddress
	serverConfig.ClusterShardLabels = CLI.ClusterShardLabels
	serverConfig.BootstrapPeers = CLI.BootstrapPeers
	serverConfig.ServiceConfig = config
	serverConfig.Logger = logger

	tp, err := tracing.InitTracing(CLI.TracingConfiguration)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to start tracing")
	}

	if tp != nil {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				logger.Warn().Err(err).Msg("failed to shutdown tracing. Spans may have been lost")
			}
		}()
	}

	var db kioradb.DB
	switch CLI.StorageBackend {
	case "boltdb":
		db, err = kioradb.NewBoltDB(CLI.StoragePath, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to create bolt db")
		}
	case "inmemory":
		db = kioradb.NewInMemoryDB()
	default:
		logger.Fatal().Msgf("unknown storage backend %s", CLI.StorageBackend)
	}

	server, err := server.NewKioraServer(serverConfig, db)
	if err != nil {
		logger.Err(err).Msg("failed to create server")
		return
	}

	// Setup a SIGINT handler, so that we can shutdown gracefully.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Info().Msg("Received signal, shutting down")
			server.Shutdown()
			break
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		logger.Err(err).Msg("failed to listen and serve")
	}

	logger.Info().Msg("Kiora Shut Down")
}
