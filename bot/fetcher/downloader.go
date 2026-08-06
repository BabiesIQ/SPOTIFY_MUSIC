/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package dl

import (
	"fmt"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"

	td "github.com/BabiesIQ/gotdbot"
)

func AcquireCachedTrack(cached *utils.CachedTrack, bot *td.Client) (string, error) {
	if cached.Platform == utils.DirectLink {
		return cached.URL, nil
	}

	if cached.Platform == utils.Telegram {
		return pullTelegramFile(cached, bot)
	}

	dlBot := bot
	if DlBot != nil {
		dlBot = DlBot
	}

	return acquireViaWrapper(cached, dlBot)
}

func acquireViaWrapper(cached *utils.CachedTrack, dlBot *td.Client) (string, error) {
	wrapper := CreateDownloaderWrapper(cached.URL)
	if !wrapper.IsValid() {
		return "", fmt.Errorf("invalid cached URL: %s", cached.URL)
	}

	track, err := wrapper.GetTrackForPlayback(cached.IsVideo)
	if err != nil {
		return "", fmt.Errorf("get track info: %w", err)
	}

	path, err := wrapper.DownloadTrack(track, cached.IsVideo)
	if err != nil {
		return "", err
	}

	if utils.TelegramMessageRegex.MatchString(path) {
		return pullFromTelegramMessage(dlBot, path)
	}

	return path, nil
}

func pullTelegramFile(cached *utils.CachedTrack, bot *td.Client) (string, error) {
	file, err := bot.GetRemoteFile(cached.TrackID, nil)
	if err != nil {
		return "", err
	}

	download, err := file.Download(bot, 0, 0, 1, &td.DownloadFileOpts{Synchronous: true})
	if err != nil {
		return "", err
	}

	return download.Local.Path, nil
}

func pullFromTelegramMessage(bot *td.Client, msgURL string) (string, error) {
	msg, err := utils.LoadMessage(bot, msgURL)
	if err != nil {
		return "", fmt.Errorf("get telegram message: %w", err)
	}

	file, err := msg.Download(bot, 1, 0, 0, true)
	if err != nil {
		return "", err
	}

	if file == nil || file.Local == nil {
		return "", fmt.Errorf("failed to download file from Telegram message")
	}

	return file.Local.Path, nil
}
