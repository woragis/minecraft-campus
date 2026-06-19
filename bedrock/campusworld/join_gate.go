package campusworld

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type playerSession struct {
	PlayerID string
	XUID     string
	Username string
	JoinedAt time.Time
	MobKills int64
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
	go gate.hudLoop()
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

	sess := playerSession{
		PlayerID: upsert.ID,
		XUID:     xuid,
		Username: name,
		JoinedAt: time.Now().UTC(),
	}
	g.sessions.Store(p, sess)
	p.Handle(&presenceHandler{gate: g, player: p, session: sess})
	g.pushHUD(ctx, p, sess)
}

func (g *JoinGate) onQuit(sess playerSession) {
	ctx := context.Background()
	sessionSeconds := int64(time.Since(sess.JoinedAt).Seconds())
	if sessionSeconds > 0 || sess.MobKills > 0 {
		if err := g.client.StatsIngest(ctx, sess.PlayerID, sessionSeconds, sess.MobKills); err != nil {
			g.log.Warn("stats ingest failed", "player", sess.Username, "err", err)
		}
	}
	if err := g.client.PresenceOffline(ctx, sess.PlayerID); err != nil {
		g.log.Warn("presence offline failed", "player", sess.Username, "err", err)
	}
}

func (g *JoinGate) addMobKill(p *player.Player) {
	value, ok := g.sessions.Load(p)
	if !ok {
		return
	}
	sess, ok := value.(playerSession)
	if !ok {
		return
	}
	sess.MobKills++
	g.sessions.Store(p, sess)
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

func (g *JoinGate) hudLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		ctx := context.Background()
		g.sessions.Range(func(key, value any) bool {
			p, ok := key.(*player.Player)
			if !ok {
				return true
			}
			sess, ok := value.(playerSession)
			if !ok {
				return true
			}
			g.pushHUD(ctx, p, sess)
			return true
		})
	}
}

func (g *JoinGate) pushHUD(ctx context.Context, p *player.Player, sess playerSession) {
	hud, err := g.client.FetchHUD(ctx, sess.PlayerID)
	if err != nil {
		g.log.Debug("hud fetch failed", "player", sess.Username, "err", err)
		return
	}
	p.SendTip(formatHUDTip(hud))
}

func formatHUDTip(hud *HUDResult) string {
	if hud == nil {
		return ""
	}
	line := hud.Username
	if hud.Status != "" {
		line += " · " + hud.Status
	}
	if hud.GuildName != "" {
		line += " · " + hud.GuildName
		if hud.GuildOnlineCount > 0 {
			line += fmt.Sprintf(" (%d online)", hud.GuildOnlineCount)
		}
	}
	return line
}

type presenceHandler struct {
	player.NopHandler
	gate    *JoinGate
	player  *player.Player
	session playerSession
}

func (h *presenceHandler) HandleQuit(p *player.Player) {
	sess := h.session
	if value, ok := h.gate.sessions.LoadAndDelete(h.player); ok {
		if latest, ok := value.(playerSession); ok {
			sess = latest
		}
	}
	h.gate.onQuit(sess)
}

func (h *presenceHandler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	if _, isPlayer := e.(*player.Player); isPlayer {
		return
	}
	if force == nil {
		return
	}
	if living, ok := e.(interface{ Health() float64 }); ok && living.Health() <= *force {
		h.gate.addMobKill(h.player)
	}
	_ = height
	_ = critical
}
