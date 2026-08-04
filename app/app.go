/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package app

import (
	"fmt"
	"log/slog"

	"github.com/BabiesIQ/gotdbot"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/store"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/config"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/player"
)

// Initialize initialises all subsystems required by the bot:
//  1. MongoDB database connection
//  2. Userbot (MTProto) assistant clients from session strings
//  3. Voice-call event handlers on the bot client
func Initialize(client *gotdbot.Client) error {
	// 1. Database
	if err := db.InitializeDatabase(); err != nil {
		return fmt.Errorf("database init failed: %w", err)
	}

	slog.Info("[App] Database initialised")

	// 2. Userbot clients
	started := 0
	for i, session := range config.SessionStrings {
		call, err := vc.Calls.StartClient(config.ApiId, config.ApiHash, session)
		if err != nil {
			slog.Warn("[App] Failed to start userbot client",
				"index", i, "error", err)
			continue
		}
		if call == nil {
			// StartClient returns nil when the account is frozen — skip silently.
			slog.Warn("[App] Userbot client is frozen, skipping", "index", i)
			continue
		}
		started++
	}

	if started == 0 {
		return fmt.Errorf("no userbot clients could be started (check STRING1–STRING10 and SESSION_TYPE)")
	}

	slog.Info("[App] Userbot clients started", "count", started)

	// 3. Voice-call handlers
	vc.Calls.RegisterHandlers(client)
	slog.Info("[App] Voice-call handlers registered")

	return nil
}
