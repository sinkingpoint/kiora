package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/internal/raft"
	"github.com/sinkingpoint/kiora/internal/server"
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

	_, err := config.LoadConfigFile(CLI.ConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	serverConfig := server.NewServerConfig()
	serverConfig.HTTPListenAddress = CLI.ListenAddress
	serverConfig.GRPCListenAddress = CLI.RaftListenURL

	config := raft.NewRaftConfig(CLI.LocalID)
	config.DataDir = CLI.RaftDataDir
	config.LocalAddress = CLI.RaftListenURL
	config.Bootstrap = CLI.RaftBootstrap
	db, err := raft.NewRaftDB(context.Background(), config, kioradb.NewInMemoryDB())
	if err != nil {
		log.Err(err).Msg("failed to initialize raft")
		return
	}

	server := server.NewKioraServer(serverConfig, db)

	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("failed to listen and serve")
	}
}
