/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package handlers

import (
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/cache"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/player"

	td "github.com/BabiesIQ/gotdbot"
)

// skipHandler handles the /skip command.
func skipHandler(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.ReplyText(c, "The bot is not streaming in the video chat.", nil)
		return nil
	}

	_ = vc.Calls.PlayNext(c, chatID)
	return nil
}
