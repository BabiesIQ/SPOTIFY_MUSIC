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

func muteController(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	if args := Args(m); args != "" {
		return td.EndGroups
	}

	chatID := m.ChatId
	if !cache.ChatCache.IsActive(chatID) {
		_, err := m.ReplyText(c, "There is no active playback in the video chat.", nil)
		return err
	}

	if _, err := vc.Calls.Mute(chatID); err != nil {
		_, err = m.ReplyText(c, fmt.Sprintf("Failed to mute the playback: %s", err.Error()), nil)
		return err
	}

	_, err := m.ReplyText(c, fmt.Sprintf("Playback has been muted by %s.", userFirstName(c, m)), &td.SendTextMessageOpts{ReplyMarkup: core.PlayerControlButtons("mute")})
	return err
}

func unmuteController(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	if args := Args(m); args != "" {
		return td.EndGroups
	}

	chatID := m.ChatId
	if !cache.ChatCache.IsActive(chatID) {
		_, err := m.ReplyText(c, "There is no active playback in the video chat.", nil)
		return err
	}

	if _, err := vc.Calls.Unmute(chatID); err != nil {
		_, err = m.ReplyText(c, fmt.Sprintf("Failed to unmute the playback: %s", err.Error()), nil)
		return err
	}

	_, err := m.ReplyText(c, fmt.Sprintf("Playback has been unmuted by %s.", userFirstName(c, m)), &td.SendTextMessageOpts{ReplyMarkup: core.PlayerControlButtons("unmute")})
	return err
}
