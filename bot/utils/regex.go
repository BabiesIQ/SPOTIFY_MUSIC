/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package utils

import "regexp"

var (
	TelegramMessageRegex = regexp.MustCompile(`^https://t\.me/(?:([a-zA-Z0-9_]{4,})|c/(\d+))/(\d+)(?:\?.*)?$`)
)
