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
	"strconv"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/cache"

	td "github.com/BabiesIQ/gotdbot"
)

func loopController(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if !cache.ChatCache.IsActive(chatID) {
		_, err := m.ReplyText(c, "There is no active playback in the video chat.", nil)
		return err
	}

	args := Args(m)
	if args == "" {
		_, err := m.ReplyText(c, "<b>Loop Control</b>\n\n<b>Usage:</b> <code>/loop [count]</code>\n0 to disable looping\n1-10 to set the number of repeats", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return err
	}

	argsInt, err := strconv.Atoi(args)
	if err != nil {
		_, _ = m.ReplyText(c, "Invalid loop value. Please provide a number between 0 and 10.", nil)
		return nil
	}

	if argsInt < 0 || argsInt > 10 {
		_, err = m.ReplyText(c, "Loop count must be between 0 and 10.", nil)
		return err
	}

	cache.ChatCache.SetLoopCount(chatID, argsInt)

	var action string
	if argsInt == 0 {
		action = "Looping has been disabled"
	} else {
		action = fmt.Sprintf("Looping has been set to %d time(s)", argsInt)
	}

	_, err = m.ReplyText(c, fmt.Sprintf("%s.\nChanged by: %s", action, userFirstName(c, m)), nil)
	return err
}
