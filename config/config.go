/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ApiId               = retrieveEnvInt32("API_ID", 0)
	ApiHash             = os.Getenv("API_HASH")
	Token               = os.Getenv("TOKEN")
	DlBotToken          = os.Getenv("DL_BOT_TOKEN")
	SessionStrings      = loadSessionStrings("STRING", 10)
	SessionType         = retrieveEnv("SESSION_TYPE", "pyrogram")
	MongoUri            = os.Getenv("MONGO_URI")
	DbName              = retrieveEnv("DB_NAME", "Anon")
	ApiUrl              = retrieveEnv("API_URL", "https://api.onegrab.fun")
	ApiKey              = os.Getenv("API_KEY")
	OwnerId             = retrieveEnvInt64("OWNER_ID", 0)
	LoggerId            = retrieveEnvInt64("LOGGER_ID", 0)
	Proxy               = os.Getenv("PROXY")
	DefaultService      = strings.ToLower(retrieveEnv("DEFAULT_SERVICE", "youtube"))
	MaxFileSize         = retrieveEnvInt64("MAX_FILE_SIZE", 500*1024*1024)
	SongDurationLimit   = retrieveEnvInt64("SONG_DURATION_LIMIT", 3600)
	DownloadsDir        = retrieveEnv("DOWNLOADS_DIR", "database")
	SupportGroup        = retrieveEnv("SUPPORT_GROUP", "https://t.me/FallenSupport")
	SupportChannel      = retrieveEnv("SUPPORT_CHANNEL", "https://t.me/FallenProjects")
	StartImg            = retrieveEnv("START_IMG", "https://i.pinimg.com/736x/0d/f4/65/0df465d1e98239ecb6283400605fc813.jpg")
	Port                = retrieveEnv("PORT", "6060")
	AutoLeave           = retrieveEnvBool("AUTO_LEAVE", false)
	EnableVideoPlayback = retrieveEnvBool("ENABLE_VPLAY", true)

	DEVS        []int64
	CookiesPath []string
	cookiesUrl  = processCookieURLs(os.Getenv("COOKIES_URL"))
)

func init() {
	devsEnv := os.Getenv("DEVS")
	if devsEnv != "" {
		devsEnv = strings.ReplaceAll(devsEnv, "\n", " ")
		devsEnv = strings.ReplaceAll(devsEnv, ",", " ")

		for _, idStr := range strings.Fields(devsEnv) {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				DEVS = append(DEVS, id)
			} else {
				slog.Info("Invalid DEV ID", "id", idStr, "error", err)
			}
		}
	}

	if OwnerId != 0 && !sliceContainsInt(DEVS, OwnerId) {
		DEVS = append(DEVS, OwnerId)
	}

	if err := validate(); err != nil {
		slog.Error("Configuration validation failed", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(DownloadsDir, 0755); err != nil {
		slog.Error("Failed to create downloads directory", "error", err)
		os.Exit(1)
	}

	if len(cookiesUrl) > 0 {
		if err := os.MkdirAll(cookiesDr, 0750); err != nil {
			slog.Error("Failed to create temp dir for cookies", "error", err)
			os.Exit(1)
		}
		go saveAllCookies(cookiesUrl)
	}
}

// retrieveEnv returns the value of an environment variable or a default value if it is not set
func retrieveEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// retrieveEnvInt64 returns the value of an environment variable as an int64 or a default value if it is not set
func retrieveEnvInt64(key string, defaultValue int64) int64 {
	if value, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
		return value
	}
	return defaultValue
}

// retrieveEnvInt32 gets environment variable as int32 with default value
func retrieveEnvInt32(key string, defaultValue int32) int32 {
	if value, err := strconv.ParseInt(os.Getenv(key), 10, 32); err == nil {
		return int32(value)
	}
	return defaultValue
}

// retrieveEnvBool gets environment variable as bool with default value
func retrieveEnvBool(key string, defaultValue bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return val
	}
	return defaultValue
}

// loadSessionStrings gets session strings from environment variable with prefix
func loadSessionStrings(prefix string, max int) []string {
	var sessions []string
	for i := 1; i <= max; i++ {
		key := fmt.Sprintf("%s%d", prefix, i)
		if session := os.Getenv(key); session != "" {
			sessions = append(sessions, session)
		}
	}

	// Also check for non-numbered version
	if session := os.Getenv(prefix); session != "" {
		sessions = append(sessions, session)
	}

	return sessions
}

// processCookieURLs processes comma-separated cookie URLs
func processCookieURLs(urls string) []string {
	if urls == "" {
		return nil
	}
	var result []string
	for _, url := range strings.Split(urls, ",") {
		url = strings.TrimSpace(url)
		if url != "" {
			result = append(result, url)
		}
	}
	return result
}

// sliceContainsInt checks if a slice sliceContains a specific int64 value
func sliceContainsInt(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// validate validates the configuration
func validate() error {
	required := []struct {
		name  string
		value string
		check func() bool
	}{
		{"API_ID", fmt.Sprintf("%d", ApiId), func() bool { return ApiId > 0 }},
		{"API_HASH", ApiHash, func() bool { return ApiHash != "" }},
		{"TOKEN", Token, func() bool { return Token != "" }},
		{"MONGO_URI", MongoUri, func() bool { return MongoUri != "" }},
		{"OWNER_ID", fmt.Sprintf("%d", OwnerId), func() bool { return OwnerId > 0 }},
	}

	var missing []string
	for _, req := range required {
		if !req.check() {
			missing = append(missing, req.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	if len(SessionStrings) == 0 {
		return fmt.Errorf("at least one session string (STRING1–10) is required")
	}

	if !isServiceValid(DefaultService) {
		DefaultService = "youtube"
		slog.Info("Invalid DEFAULT_SERVICE, defaulting to 'youtube'", "Service", DefaultService)
	}

	return nil
}

// isServiceValid checks if the service is valid
func isServiceValid(service string) bool {
	validServices := map[string]bool{
		"youtube": true,
		"spotify": true,
	}
	return validServices[strings.ToLower(service)]
}
