package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/woragis/minecraft-campus-bedrock/campusworld"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	log := slog.Default()
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(log)
	if err != nil {
		log.Error("config", "err", err)
		os.Exit(1)
	}

	apiURL := envOr("CAMPUS_API_URL", "http://127.0.0.1:8080")
	apiKey := strings.TrimSpace(os.Getenv("PLUGIN_API_KEY"))
	if apiKey == "" {
		log.Error("PLUGIN_API_KEY is required")
		os.Exit(1)
	}
	serverSlug := envOr("SERVER_SLUG", "bedrock")

	client := campusworld.NewClient(apiURL, apiKey, serverSlug)
	gate := campusworld.NewJoinGate(client, log)

	srv := conf.New()
	srv.CloseOnProgramEnd()
	srv.Listen()

	log.Info("campusworld bedrock listening", "api", apiURL, "slug", serverSlug)

	for p := range srv.Accept() {
		gate.Handle(p)
	}
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	if name := strings.TrimSpace(os.Getenv("SERVER_NAME")); name != "" {
		c.Server.Name = name
	}
	return c.Config(log)
}
