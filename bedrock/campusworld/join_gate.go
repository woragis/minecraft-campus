package campusworld

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/player"
)

type playerSession struct {
	PlayerID string
	XUID     string
	Username string
}

type JoinGate struct {
	client   *Client
	log      *slog.Logger
	sessions sync.Map
}

func NewJoinGate(client *Client, log *slog.Logger) *JoinGate {
	if log == nil {
		log = slog.Default()
	}
	gate := &JoinGate{client: client, log: log}
	go gate.heartbeatLoop()
	return gate
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

	upsert, err := g.client.UpsertBedrockPlayer(ctx, xuid, name)
	if err != nil {
		g.log.Warn("player upsert failed", "player", name, "xuid", xuid, "err", err)
		return
	}
	if upsert == nil || upsert.ID == "" {
		return
	}

	sess := playerSession{PlayerID: upsert.ID, XUID: xuid, Username: name}
	g.sessions.Store(p, sess)
	p.Handle(&presenceHandler{gate: g, player: p, session: sess})
}

func (g *JoinGate) onQuit(p *player.Player, sess playerSession) {
	g.sessions.Delete(p)
	ctx := context.Background()
	if err := g.client.PresenceOffline(ctx, sess.PlayerID); err != nil {
		g.log.Warn("presence offline failed", "player", sess.Username, "err", err)
	}
}

func (g *JoinGate) heartbeatLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		ctx := context.Background()
		g.sessions.Range(func(key, value any) bool {
			sess, ok := value.(playerSession)
			if !ok {
				return true
			}
			if err := g.client.PresenceHeartbeat(ctx, sess.PlayerID); err != nil {
				g.log.Debug("presence heartbeat failed", "player", sess.Username, "err", err)
			}
			return true
		})
	}
}

type presenceHandler struct {
	player.NopHandler
	gate    *JoinGate
	player  *player.Player
	session playerSession
}

func (h *presenceHandler) HandleQuit(p *player.Player) {
	h.gate.onQuit(h.player, h.session)
}
