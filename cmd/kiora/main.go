package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/internal/server"
	"github.com/sinkingpoint/kiora/internal/tracing"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

var CLI struct {
	ListenAddress string `name:"web.listen-url" help:"the address to listen on" default:"localhost:4278"`
	ConfigFile    string `name:"config.file" short:"c" help:"the config file to load config from" default:"./kiora.dot"`
	RaftDataDir   string `name:"raft.data-dir" help:"the directory to put database state in" default:"./kiora/data"`
	RaftBootstrap bool   `name:"raft.bootstrap" help:"If set, bootstrap a new raft cluster"`
	LocalID       string `name:"raft.local-id" help:"the name of this node in the raft cluster" default:""`
	RaftListenURL string `name:"raft.listen-url" help:"the address for the raft node to listen on" default:"localhost:4279"`
}

func main() {
	kong.Parse(&CLI)

	if CLI.LocalID == "" {
		var err error
		CLI.LocalID, err = os.Hostname()
		if err != nil {
			log.Fatal().Err(err).Msg("no local id set, and failed to get hostname")
		}
	}

	conf, err := config.LoadConfigFile(CLI.ConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	serverConfig := server.NewServerConfig()
	serverConfig.HTTPListenAddress = CLI.ListenAddress
	serverConfig.GRPCListenAddress = CLI.RaftListenURL

	serverConfig.RaftConfig.LocalID = CLI.LocalID
	serverConfig.RaftConfig.DataDir = CLI.RaftDataDir
	serverConfig.RaftConfig.LocalAddress = CLI.RaftListenURL
	serverConfig.RaftConfig.Bootstrap = CLI.RaftBootstrap
	serverConfig.NotifyConfig = conf

	tracingConfig := tracing.DefaultTracingConfiguration()
	tp, err := tracing.InitTracing(tracingConfig)
	if err != nil {
		log.Warn().Err(err).Msg("failed to start tracing")
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Warn().Err(err).Msg("failed to shutdown tracing. Spans may have been lost")
		}
	}()

	server, err := server.NewKioraServer(serverConfig, kioradb.NewInMemoryDB())
	if err != nil {
		log.Err(err).Msg("failed to create server")
		return
	}

	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("failed to listen and serve")
	}
}
