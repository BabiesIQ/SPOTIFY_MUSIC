/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package vc

import (
	"log/slog"
	"regexp"
	"sync"
	"time"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/cache"

	td "github.com/BabiesIQ/gotdbot"
	tg "github.com/BabiesIQ/gogram/telegram"
)

var logger = slog.Default()
var urlRegex = regexp.MustCompile(`^https?://`)

// TelegramCalls manages the state and operations for voice calls, including userbots and the main bot client.
type TelegramCalls struct {
	mu          sync.RWMutex
	assistants  map[int]*Assistant
	clients     map[int]*tg.Client
	statusCache *cache.Cache[td.ChatMemberStatus]
	inviteCache *cache.Cache[string]

	leavingMu sync.Mutex
	leaving   map[int]bool
}

var (
	instance *TelegramCalls
	once     sync.Once
)

// activeCalls returns the singleton instance of the TelegramCalls manager, ensuring that only one instance is created.
func activeCalls() *TelegramCalls {
	once.Do(func() {
		instance = &TelegramCalls{
			assistants:  make(map[int]*Assistant),
			clients:     make(map[int]*tg.Client),
			statusCache: cache.NewCache[td.ChatMemberStatus](2 * time.Hour),
			inviteCache: cache.NewCache[string](2 * time.Hour),
			leaving:     make(map[int]bool),
		}
	})
	return instance
}

// Calls is the singleton instance of TelegramCalls, initialized lazily.
var Calls = activeCalls()
