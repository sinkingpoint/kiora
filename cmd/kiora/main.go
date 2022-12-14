package main

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/internal/raft"
	"github.com/sinkingpoint/kiora/internal/server"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

var CLI struct {
	ListenAddress string `name:"web.listen-url" help:"the address to listen on" default:"localhost:4278"`
	ConfigFile    string `name:"config.file" help:"the config file to load config from" default:"./kiora.toml"`
}

func main() {
	kong.Parse(&CLI)

	_, err := config.LoadConfigFile(CLI.ConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	serverConfig := server.NewServerConfig()
	serverConfig.ListenAddress = CLI.ListenAddress

	config := raft.NewRaftConfig("local")
	stateMachine := raft.NewAlertTracker(kioradb.NewInMemoryDB())
	_, err = raft.NewRaft(context.Background(), config, stateMachine)
	if err != nil {
		log.Err(err).Msg("failed to initialize raft")
		return
	}

	server := server.NewKioraServer(serverConfig, kioradb.NewInMemoryDB())
	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("failed to listen and serve")
	}
}
