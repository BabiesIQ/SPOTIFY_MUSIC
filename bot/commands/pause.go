/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package handlers

import (
	"fmt"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/keyboard"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/cache"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/player"

	td "github.com/BabiesIQ/gotdbot"
)

func pauseHandler(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.ReplyText(c, "There is no active playback in the video chat.", nil)
		return nil
	}

	if _, err := vc.Calls.Pause(chatID); err != nil {
		_, _ = m.ReplyText(c, fmt.Sprintf("Failed to pause the playback: %s", err.Error()), nil)
		return nil
	}

	_, err := m.ReplyText(c, fmt.Sprintf("Playback has been paused by %s.", userFirstName(c, m)), &td.SendTextMessageOpts{ReplyMarkup: core.PlayerControlButtons("pause")})
	return err
}

func resumeHandler(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if chatID > 0 {
		_, _ = m.ReplyText(c, "This command can only be used in a supergroup.", nil)
		return nil
	}

	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.ReplyText(c, "There is no active playback in the video chat.", nil)
		return nil
	}

	if _, err := vc.Calls.Resume(chatID); err != nil {
		_, _ = m.ReplyText(c, fmt.Sprintf("Failed to resume the playback: %s", err.Error()), nil)
		return nil
	}

	_, err := m.ReplyText(c, fmt.Sprintf("Playback has been resumed by %s.", userFirstName(c, m)), &td.SendTextMessageOpts{ReplyMarkup: core.PlayerControlButtons("resume")})
	return err
}
