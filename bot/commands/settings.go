/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package handlers

import (
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"
	"fmt"
	"strings"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/keyboard"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/cache"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/store"

	td "github.com/BabiesIQ/gotdbot"
)

func settingsHandler(c *td.Client, m *td.Message) error {
	if m.IsPrivate() {
		return nil
	}

	if !adminModeState(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId
	admins, err := cache.LoadAdmins(c, chatID, false)
	if err != nil {
		return err
	}

	// Check if user is admin
	var hasAdminRights bool
	for _, admin := range admins {
		if SenderID(admin.MemberId) == m.SenderID() {
			hasAdminRights = true
			break
		}
	}

	if !hasAdminRights {
		return nil
	}

	// Get current settings
	getPlayMode := db.Instance.GetPlayMode(chatID)
	playModeStr := utils.Everyone
	if getPlayMode {
		playModeStr = utils.Admins
	}
	getAdminMode := db.Instance.GetAdminMode(chatID)
	cmdDelete := db.Instance.GetCmdDelete(chatID)
	language, _ := db.Instance.GetLanguage(chatID)

	chat, err := m.GetChat(c)
	if err != nil {
		c.Logger.Warn("Failed to get chat", "error", err)
		return nil
	}

	text := fmt.Sprintf("<u><b>%s settings</b></u>\n\nClick the buttons below to change this chat's current settings.",
		chat.Title)

	_, err = m.ReplyText(c, text, &td.SendTextMessageOpts{ReplyMarkup: core.SettingsKeyboard(playModeStr, getAdminMode, cmdDelete, language), ParseMode: td.ParseModeHTML})
	return err
}

func settingsCallbackHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	if cb.IsPrivate() {
		return nil
	}

	chatID := cb.ChatId

	// Check admin permissions
	admins, err := cache.LoadAdmins(c, chatID, false)
	if err != nil {
		return err
	}

	var hasPerms bool
	for _, admin := range admins {
		if SenderID(admin.MemberId) == cb.SenderUserId {
			rights, _ := cache.LoadRights(c, chatID, cb.SenderUserId, false)
			hasPerms = (rights != nil && rights.CanManageVideoChats) || admin.Status == td.ChatMemberStatusCreator{}
			break
		}
	}

	if !hasPerms {
		err = cb.Answer(c, 0, true, "You don't have permission to change settings.", "")
		return err
	}

	// Process the callback data
	data := cb.DataString()
	if data == "settings_main" {
		return cb.Answer(c, 0, false, "Update your chat settings", "")
	}

	parts := strings.Split(data, "_")
	if len(parts) < 2 {
		return nil
	}

	settingType := parts[1]

	switch settingType {
	case "delete":
		cmdDelete := db.Instance.GetCmdDelete(chatID)
		_ = db.Instance.SetCmdDelete(chatID, !cmdDelete)
	case "play":
		getPlayMode := db.Instance.GetPlayMode(chatID)
		_ = db.Instance.SetPlayMode(chatID, !getPlayMode)
	case "admin":
		getAdminMode := db.Instance.GetAdminMode(chatID)
		newMode := utils.Everyone
		if getAdminMode == utils.Everyone {
			newMode = utils.Admins
		}
		_ = db.Instance.SetAdminMode(chatID, newMode)
	case "lang":
		return cb.Answer(c, 0, true, "Language selection is not yet implemented via this menu.", "")
	default:
		return cb.Answer(c, 0, true, "Unknown setting", "")
	}

	getPlayMode := db.Instance.GetPlayMode(chatID)
	playModeStr := utils.Everyone
	if getPlayMode {
		playModeStr = utils.Admins
	}
	getAdminMode := db.Instance.GetAdminMode(chatID)
	cmdDelete := db.Instance.GetCmdDelete(chatID)
	language, _ := db.Instance.GetLanguage(chatID)

	chat, err := c.GetChat(chatID)
	if err != nil {
		c.Logger.Warn("Failed to get chat", "error", err)
		return nil
	}

	text := fmt.Sprintf("<u><b>%s settings</b></u>\n\nClick the buttons below to change this chat's current settings.",
		chat.Title)

	_, err = cb.EditMessageText(c, text, &td.EditTextMessageOpts{ReplyMarkup: core.SettingsKeyboard(playModeStr, getAdminMode, cmdDelete, language), ParseMode: td.ParseModeHTML})
	if err != nil {
		return err
	}

	_ = cb.Answer(c, 0, false, "Settings updated", "")
	return nil
}
