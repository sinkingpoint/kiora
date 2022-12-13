package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/cmd/kiora/config"
	"github.com/sinkingpoint/kiora/internal/server"
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

	server := server.NewKioraServer(serverConfig)
	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("failed to listen and serve")
	}
}
