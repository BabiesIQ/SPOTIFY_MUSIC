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
	"fmt"

	td "github.com/BabiesIQ/gotdbot"
)

func autoplayController(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if cache.ChatCache.GetPlayingTrack(chatID) == nil {
		_, err := m.ReplyText(c, "Bot is not streaming in the video chat.", nil)
		return err
	}

	state := cache.ChatCache.GetAutoplay(chatID)
	text := "<b>Autoplay Control</b>\n\nWhen autoplay is enabled, the bot will automatically play recommended songs from YouTube when the queue is empty."
	button := buildAutoplayButton(state)

	_, err := m.ReplyText(c, text, &td.SendTextMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: button,
	})
	return err
}

func callbackAutoplayHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	if !adminModeCallback(c, cb) {
		return nil
	}

	chatID := cb.ChatId
	if cache.ChatCache.GetPlayingTrack(chatID) == nil {
		_ = cb.Answer(c, 0, true, "Bot is not streaming in the video chat.", "")
		return nil
	}

	state := cache.ChatCache.GetAutoplay(chatID)
	newState := !state
	cache.ChatCache.SetAutoplay(chatID, newState)

	text := "<b>Autoplay Control</b>\n\nWhen autoplay is enabled, the bot will automatically play recommended songs from YouTube when the queue is empty."
	button := buildAutoplayButton(newState)

	_, err := cb.EditMessageText(c, text, &td.EditTextMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: button,
	})
	if err != nil {
		c.Logger.Warn("Failed to edit autoplay message", "error", err)
	}

	var status string
	if newState {
		status = "enabled"
	} else {
		status = "disabled"
	}
	_ = cb.Answer(c, 0, false, fmt.Sprintf("Autoplay has been %s.", status), "")

	return nil
}

func buildAutoplayButton(state bool) *td.ReplyMarkupInlineKeyboard {
	var text string
	if state {
		text = "Autoplay: ON | ✅"
	} else {
		text = "Autoplay: OFF | ❌"
	}

	return &td.ReplyMarkupInlineKeyboard{
		Rows: [][]td.InlineKeyboardButton{
			{
				{
					Text: text,
					Type: &td.InlineKeyboardButtonTypeCallback{
						Data: []byte("autoplay_toggle"),
					},
				},
			},
		},
	}
}
