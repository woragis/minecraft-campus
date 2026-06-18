package campusworld

import (
	"context"
	"log/slog"

	"github.com/df-mc/dragonfly/server/player"
)

type JoinGate struct {
	client *Client
	log    *slog.Logger
}

func NewJoinGate(client *Client, log *slog.Logger) *JoinGate {
	if log == nil {
		log = slog.Default()
	}
	return &JoinGate{client: client, log: log}
}

func (g *JoinGate) Handle(p *player.Player) {
	xuid := p.XUID()
	name := p.Name()
	if xuid == "" {
		p.Disconnect("XBOX Live sign-in is required to join CampusWorld.")
		return
	}

	ctx := context.Background()
	result, err := g.client.CheckBedrockWhitelist(ctx, xuid, name)
	if err != nil {
		g.log.Error("whitelist check failed", "player", name, "xuid", xuid, "err", err)
		p.Disconnect("CampusWorld is unavailable. Try again in a moment.")
		return
	}
	if !result.Allowed {
		p.Disconnect(KickMessage(result.Reason))
		return
	}

	if err := g.client.UpsertBedrockPlayer(ctx, xuid, name); err != nil {
		g.log.Warn("player upsert failed", "player", name, "xuid", xuid, "err", err)
	}
}
