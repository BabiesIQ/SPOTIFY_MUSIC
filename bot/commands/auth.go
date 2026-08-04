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
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/store"
	"fmt"

	td "github.com/BabiesIQ/gotdbot"
)

func listAuthHandler(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	if m.IsPrivate() {
		return nil
	}

	chatID := m.ChatId

	authUser := db.Instance.GetAuthUsers(chatID)
	if authUser == nil || len(authUser) == 0 {
		_, _ = m.ReplyText(c, "No authorized users found.", nil)
		return nil
	}

	text := "<b>Authorized Users</b>\n\n"
	for _, uid := range authUser {
		text += fmt.Sprintf("• <a href=\"tg://user?id=%d\">%d</a>\n", uid, uid)
	}

	_, _ = m.ReplyText(c, text, replyOpts)
	return td.EndGroups
}

func authAddHandler(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	if m.IsPrivate() {
		return td.EndGroups
	}

	chatID := m.ChatId

	UserStatus, err := cache.LoadUserAdmin(c, chatID, m.SenderID(), false)
	if err != nil {
		c.Logger.Warn("LoadUserAdmin error", "error", err)
		_, _ = m.ReplyText(c, "Unable to verify administrator status.", nil)
		return td.EndGroups
	}

	switch UserStatus.Status.(type) {
	case *td.ChatMemberStatusCreator, *td.ChatMemberStatusAdministrator:
	default:
		_, _ = m.ReplyText(c, "You must be an administrator to use this command.", nil)
		return td.EndGroups
	}

	userID, err := resolveTargetUserID(c, m)
	if err != nil {
		_, _ = m.ReplyText(c, err.Error(), nil)
		return nil
	}

	if db.Instance.IsAuthUser(chatID, userID) {
		_, _ = m.ReplyText(c, "This user is already authorized.", nil)
		return nil
	}

	if err = db.Instance.AddAuthUser(chatID, userID); err != nil {
		c.Logger.Error("Failed to add authorized user", "error", err)
		_, _ = m.ReplyText(c, "Failed to authorize the user.", nil)
		return nil
	}

	_, err = m.ReplyText(c, fmt.Sprintf("User %d has been authorized.", userID), nil)
	return err
}

func removeAuthHandler(c *td.Client, m *td.Message) error {
	if !adminModeState(c, m) {
		return td.EndGroups
	}

	if m.IsPrivate() {
		return td.EndGroups
	}

	chatID := m.ChatId

	UserStatus, err := cache.LoadUserAdmin(c, chatID, m.SenderID(), false)
	if err != nil {
		c.Logger.Warn("LoadUserAdmin error", "error", err)
		_, _ = m.ReplyText(c, "Unable to verify administrator status.", nil)
		return td.EndGroups
	}

	switch UserStatus.Status.(type) {
	case *td.ChatMemberStatusCreator, *td.ChatMemberStatusAdministrator:
	default:
		_, _ = m.ReplyText(c, "You must be an administrator to use this command.", nil)
		return td.EndGroups
	}

	userID, err := resolveTargetUserID(c, m)
	if err != nil {
		_, _ = m.ReplyText(c, err.Error(), nil)
		return nil
	}

	if !db.Instance.IsAuthUser(chatID, userID) {
		_, _ = m.ReplyText(c, "This user is not authorized.", nil)
		return nil
	}

	if err := db.Instance.RemoveAuthUser(chatID, userID); err != nil {
		c.Logger.Error("Failed to remove authorized user", "error", err)
		_, _ = m.ReplyText(c, "Failed to remove authorized user.", nil)
		return nil
	}

	_, err = m.ReplyText(c, fmt.Sprintf("User %d has been removed from the authorized list.", userID), nil)
	return err
}
