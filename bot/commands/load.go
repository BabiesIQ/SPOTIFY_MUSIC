/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package handlers

import (
	"time"

	"github.com/BabiesIQ/gotdbot"
	"github.com/BabiesIQ/gotdbot/filters/callbackquery"
)

var startTime = time.Now()

// RegisterModules loads all the handlers.
// It takes a telegram gotdbot.Client as input.
func RegisterModules(c *gotdbot.Client) {
	c.OnCommand("reload", reloadAdminCacheHandler)
	c.OnCommand("authList", listAuthHandler)
	c.OnCommand("auths", listAuthHandler)
	c.OnCommand("auth", authAddHandler)
	c.OnCommand("addAuth", authAddHandler)
	c.OnCommand("removeAuth", removeAuthHandler)
	c.OnCommand("rmAuth", removeAuthHandler)
	c.OnCommand("broadcast", broadcastController)
	c.OnCommand("gCast", broadcastController)
	c.OnCommand("stop_gcast", broadcastCancelHandler)
	c.OnCommand("stop_broadcast", broadcastCancelHandler)
	c.OnCommand("av", vcActiveHandler)
	c.OnCommand("active_vc", vcActiveHandler)
	c.OnCommand("clearass", assistantsClearHandler)
	c.OnCommand("clearAssistants", assistantsClearHandler)
	c.OnCommand("leaveAll", allLeaveHandler)
	c.OnCommand("logger", logController)
	c.OnCommand("privacy", privacyHandler)
	c.OnCommand("autoplay", autoplayController)
	c.OnCommand("loop", loopController)
	c.OnCommand("pause", pauseHandler)
	c.OnCommand("resume", resumeHandler)
	c.OnCommand("cplist", playlistCreateHandler)
	c.OnCommand("createplaylist", playlistCreateHandler)
	c.OnCommand("deleteplaylist", playlistDeleteHandler)
	c.OnCommand("queue", queueHandler)
	c.OnCommand("seek", seekHandler)
	c.OnCommand("sh", shellCommand)
	c.OnCommand("skip", skipHandler)
	c.OnCommand("stop", stopHandler)
	c.OnCommand("end", stopHandler)
	c.OnCommand("start", startHandler)
	c.OnCommand("help", startHandler)
	c.OnCommand("ping", pingHandler)
	c.OnCommand("play", playHandler)
	c.OnCommand("p", playHandler)
	c.OnCommand("fplay", playFileHandler)
	c.OnCommand("fp", playFileHandler)
	c.OnCommand("vplay", playVideoHandler)
	c.OnCommand("v", playVideoHandler)
	c.OnCommand("fvplay", playVideoFileHandler)
	c.OnCommand("fvp", playVideoFileHandler)
	c.OnCommand("remove", removeHandler)
	c.OnCommand("mute", muteController)
	c.OnCommand("unmute", unmuteController)
	c.OnCommand("settings", settingsHandler)
	c.OnCommand("addtoplaylist", playlistAddHandler)
	c.OnCommand("addtoplist", playlistAddHandler)
	c.OnCommand("removefromplaylist", removeFromPlaylistHandler)
	c.OnCommand("rmplist", removeFromPlaylistHandler)
	c.OnCommand("plistinfo", playlistInfoHandler)
	c.OnCommand("playlistinfo", playlistInfoHandler)
	c.OnCommand("myplaylists", userPlaylistsHandler)
	c.OnCommand("myplist", userPlaylistsHandler)
	c.OnCommand("stats", statsHandler)

	c.OnUpdateNewCallbackQuery(callbackHelpHandler, callbackquery.Prefix("help_"))
	c.OnUpdateNewCallbackQuery(playCallbackHandler, callbackquery.Prefix("play_"))
	c.OnUpdateNewCallbackQuery(playVcHandler, callbackquery.Prefix("vcplay_"))
	c.OnUpdateNewCallbackQuery(settingsCallbackHandler, callbackquery.Prefix("settings_"))
	c.OnUpdateNewCallbackQuery(callbackAutoplayHandler, callbackquery.Equal("autoplay_toggle"))

	c.OnUpdateChatMember(processParticipant, nil)
	c.OnUpdateNewMessage(processVoiceChatMessage, nil)

	c.Logger.Debug("Handlers loaded successfully")
}
