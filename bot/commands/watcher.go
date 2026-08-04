/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package handlers

import (
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/keyboard"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/cache"
	"fmt"
	"time"

	td "github.com/BabiesIQ/gotdbot"
)

func processVoiceChatMessage(c *td.Client, update *td.UpdateNewMessage) error {
	m := update.Message
	chatID := m.ChatId

	if m.IsGroup() {
		text := fmt.Sprintf(
			"This chat (%d) is not a supergroup yet.\n<b>⚠️ Please convert this chat to a supergroup and add me as admin.</b>\n\nIf you don't know how to convert, use this guide:\n🔗 https://te.legra.ph/How-to-Convert-a-Group-to-a-Supergroup-01-02\n\nIf you have any questions, join our support group:",
			chatID,
		)

		_, _ = c.SendTextMessage(chatID, text, &td.SendTextMessageOpts{
			ReplyMarkup:           core.JoinButtonMarkup(c.Me.Usernames.EditableUsername),
			DisableWebPagePreview: true,
			ParseMode:             "HTML",
		})

		time.Sleep(1 * time.Second)
		_ = c.LeaveChat(chatID)
		return nil
	}

	go storeChatToDB(chatID)

	if m.Content == nil {
		return nil
	}
	var message string
	switch m.Content.(type) {
	case *td.MessageVideoChatStarted:
		cache.ChatCache.ClearChat(chatID)
		message = "🎙️ Video chat started!\nUse /play <song name> to play music."
	case *td.MessageVideoChatEnded:
		cache.ChatCache.ClearChat(chatID)
		message = "🎧 Video chat ended!\nAll queues cleared."
	default:
		return nil
	}

	_, _ = c.SendTextMessage(chatID, message, nil)
	return td.EndGroups
}
