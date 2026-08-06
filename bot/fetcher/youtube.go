/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package dl

import (
	"time"

	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/config"
	"log/slog"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// youTubeData provides an interface for fetching track and playlist information from YouTube.
type youTubeData struct {
	Query    string
	ApiUrl   string
	APIKey   string
	Patterns map[string]*regexp.Regexp
}

var youtubePatterns = map[string]*regexp.Regexp{
	"youtube":   regexp.MustCompile(`(?i)^(?:https?://)?(?:www\.)?youtube\.com/.*`),
	"youtu_be":  regexp.MustCompile(`(?i)^(?:https?://)?(?:www\.)?youtu\.be/.*`),
	"yt_music":  regexp.MustCompile(`(?i)^(?:https?://)?music\.youtube\.com/.*`),
	"yt_shorts": regexp.MustCompile(`(?i)^(?:https?://)?(?:www\.)?youtube\.com/shorts/.*`),
}

var (
	youtubeSearch = searchYouTube
	ytDlpDownload = func(y *youTubeData, videoID string, video bool) (string, error) {
		return y.downloadWithYtDlp(videoID, video)
	}
)

// createYouTubeData initializes a youTubeData instance with pre-compiled regex patterns and a cleaned query.
func createYouTubeData(query string) *youTubeData {
	return &youTubeData{
		Query:    strings.TrimSpace(query),
		ApiUrl:   strings.TrimRight(config.ApiUrl, "/"),
		APIKey:   config.ApiKey,
		Patterns: youtubePatterns,
	}
}

func (y *youTubeData) isValid() bool {
	if y.Query == "" {
		slog.Info("The query or patterns are empty.")
		return false
	}

	for _, pattern := range y.Patterns {
		if pattern.MatchString(y.Query) {
			return true
		}
	}

	// A plain text query is also a YouTube request. URLs that belong to
	// another supported service are intentionally left for that service.
	if strings.Contains(y.Query, "://") || strings.HasPrefix(strings.ToLower(y.Query), "www.") {
		return false
	}
	for _, pattern := range apiPatterns {
		if pattern.MatchString(y.Query) {
			return false
		}
	}
	return true
}

func (y *youTubeData) getInfo() (utils.PlatformTracks, error) {
	if !y.isValid() {
		return utils.PlatformTracks{}, errors.New("the provided URL is invalid or the platform is not supported")
	}

	if tracks, err := y.getInfoWithBabyAPI(false); err == nil {
		return tracks, nil
	}

	return y.getInfoFromYouTube()
}

func (y *youTubeData) getInfoFromYouTube() (utils.PlatformTracks, error) {
	if !y.isValid() {
		return utils.PlatformTracks{}, errors.New("the provided URL is invalid or the platform is not supported")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	y.Query = normalizeYouTubeURL(y.Query)
	videoID := parseVideoID(y.Query)
	playlistID := parsePlaylistID(y.Query)

	switch {
	case playlistID != "":
		if strings.HasPrefix(playlistID, "RD") {
			return FetchYouTubeMixPlaylist(ctx, playlistID)
		}
		return fetchYouTubePlaylist(ctx, playlistID)

	case videoID != "":
		for _, query := range []string{videoID, y.Query} {
			tracks, err := youtubeSearch(query, 10)
			if err != nil {
				continue
			}

			for _, track := range tracks {
				if track.Id == videoID {
					return utils.PlatformTracks{Results: []utils.MusicTrack{track}}, nil
				}
			}
		}

		if title, err := fetchYouTubeTitleFromOEmbed(videoID); err == nil && title != "" {
			tracks, err := youtubeSearch(title, 10)
			if err == nil {
				for _, track := range tracks {
					if track.Id == videoID {
						return utils.PlatformTracks{Results: []utils.MusicTrack{track}}, nil
					}
				}
			}
		}

		slog.Warn("Video ID was extracted but no matching track was found in search results", "video_id", videoID)
		return fetchYouTubeVideo(ctx, videoID)
	}

	tracks, err := youtubeSearch(y.Query, 5)
	if err != nil {
		return utils.PlatformTracks{}, fmt.Errorf("YouTube search fallback failed: %w", err)
	}
	if len(tracks) == 0 {
		return utils.PlatformTracks{}, errors.New("no video or playlist results were found")
	}
	return utils.PlatformTracks{Results: tracks}, nil
}

func (y *youTubeData) search() (utils.PlatformTracks, error) {
	return y.searchForPlayback(false)
}

func (y *youTubeData) searchForPlayback(video bool) (utils.PlatformTracks, error) {
	if tracks, err := y.searchWithBabyAPI(video); err == nil {
		return tracks, nil
	} else {
		slog.Warn("BabyAPI search failed; falling back to YouTube search", "query", y.Query, "video", video, "error", err)
	}

	tracks, err := youtubeSearch(y.Query, 5)
	if err != nil {
		return utils.PlatformTracks{}, err
	}

	if len(tracks) == 0 {
		return utils.PlatformTracks{}, errors.New("no video results were found")
	}

	return utils.PlatformTracks{Results: tracks}, nil
}

func (y *youTubeData) getTrack() (utils.TrackInfo, error) {
	if y.Query == "" {
		return utils.TrackInfo{}, errors.New("the query is empty")
	}

	if !y.isValid() {
		return utils.TrackInfo{}, errors.New("the provided URL is invalid or the platform is not supported")
	}

	if trackInfo, err := y.getTrackWithBabyAPI(false); err == nil {
		return trackInfo, nil
	} else {
		slog.Warn("BabyAPI track lookup failed; falling back to YouTube", "query", y.Query, "error", err)
	}

	getInfo, err := y.getInfoFromYouTube()
	if err != nil {
		return utils.TrackInfo{}, err
	}
	if len(getInfo.Results) == 0 {
		return utils.TrackInfo{}, errors.New("no video results were found")
	}

	track := getInfo.Results[0]
	trackInfo := utils.TrackInfo{
		Id:       track.Id,
		URL:      track.Url,
		Platform: utils.YouTube,
	}

	return trackInfo, nil
}

func (y *youTubeData) getTrackForPlayback(video bool) (utils.TrackInfo, error) {
	if y.Query == "" {
		return utils.TrackInfo{}, errors.New("the query is empty")
	}
	if !y.isValid() {
		return utils.TrackInfo{}, errors.New("the provided URL is invalid or the platform is not supported")
	}

	if trackInfo, err := y.getTrackWithBabyAPI(video); err == nil {
		return trackInfo, nil
	} else {
		slog.Warn("BabyAPI track lookup failed; falling back to YouTube", "query", y.Query, "video", video, "error", err)
	}

	getInfo, err := y.getInfoFromYouTube()
	if err != nil {
		return utils.TrackInfo{}, err
	}
	if len(getInfo.Results) == 0 {
		return utils.TrackInfo{}, errors.New("no video results were found")
	}

	track := getInfo.Results[0]
	return utils.TrackInfo{
		Id:       track.Id,
		URL:      track.Url,
		Platform: utils.YouTube,
	}, nil
}

// downloadTrack handles the download of a track from YouTube.
func (y *youTubeData) downloadTrack(info utils.TrackInfo, video bool) (string, error) {
	if !video && info.CdnURL != "" {
		return info.CdnURL, nil
	}

	if filePath, err := y.downloadWithBabyAPI(info.Id, video); err == nil {
		return filePath, nil
	} else {
		slog.Warn("BabyAPI download failed; falling back to yt-dlp", "video_id", info.Id, "video", video, "error", err)
	}

	return ytDlpDownload(y, info.Id, video)
}

// buildYtdlpParams constructs the command-line parameters for yt-dlp to download media.
func (y *youTubeData) buildYtdlpParams(videoID string, video bool) ([]string, string) {
	outputTemplate := filepath.Join(config.DownloadsDir, "%(id)s.%(ext)s")
	var cookieFile string

	params := []string{
		"yt-dlp",
		"--no-warnings",
		"--quiet",
		"--geo-bypass",
		"--retries", "2",
		"--continue",
		"--no-part",
		"--concurrent-fragments", "3",
		"--socket-timeout", "10",
		"--throttled-rate", "100K",
		"--retry-sleep", "1",
		"--no-write-thumbnail",
		"--no-write-info-json",
		"--no-embed-metadata",
		"--no-embed-chapters",
		"--no-embed-subs",
		"--extractor-args", "youtube:player_js_version=actual",
		"-o", outputTemplate,
	}

	if video {
		formatSelector := "bestvideo[height<=720]+bestaudio/best[height<=720]"
		params = append(params, "-f", formatSelector, "--merge-output-format", "mp4")
	} else {
		params = append(params, "-f", "bestaudio[ext=m4a]/bestaudio")
	}

	cookieFile = y.getCookieFile()
	if cookieFile != "" {
		params = append(params, "--cookies", cookieFile)
	} else if config.Proxy != "" {
		params = append(params, "--proxy", config.Proxy)
	}

	videoURL := "https://www.youtube.com/watch?v=" + videoID
	params = append(params, videoURL, "--print", "after_move:filepath")

	return params, cookieFile
}

// downloadWithYtDlp downloads media from YouTube using the yt-dlp command-line tool.
func (y *youTubeData) downloadWithYtDlp(videoID string, video bool) (string, error) {
	if videoID == "" {
		return "", errors.New("videoID is empty")
	}

	ytdlpParams, cookieFile := y.buildYtdlpParams(videoID, video)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, ytdlpParams[0], ytdlpParams[1:]...)

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			stderr := string(exitErr.Stderr)
			if cookieFile != "" && strings.Contains(stderr, "Sign in to confirm you're not a bot") {
				_ = os.Remove(cookieFile)
			}
			return "", fmt.Errorf("yt-dlp failed with exit code %d: %s", exitErr.ExitCode(), stderr)
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("yt-dlp timed out for video ID: %s", videoID)
		}

		return "", fmt.Errorf("an unexpected error occurred while downloading %s: %w", videoID, err)
	}

	downloadedPathStr := strings.TrimSpace(string(output))
	if downloadedPathStr == "" {
		return "", fmt.Errorf("no output path was returned for %s", videoID)
	}

	if _, err := os.Stat(downloadedPathStr); os.IsNotExist(err) {
		return "", fmt.Errorf("the file was not found at the reported path: %s", downloadedPathStr)
	}

	return downloadedPathStr, nil
}

// getCookieFile retrieves the path to a cookie file from the configured list.
func (y *youTubeData) getCookieFile() string {
	cookiesPath := config.CookiesPath
	if len(cookiesPath) == 0 {
		return ""
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(cookiesPath))))
	if err != nil {
		slog.Info("Could not generate a random number", "error", err)
		return cookiesPath[0]
	}

	return cookiesPath[n.Int64()]
}

// downloadWithApi downloads a track using the external API.
func (y *youTubeData) downloadWithApi(videoID string, _ bool) (string, error) {
	videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	api := createApiData(videoUrl)
	track, err := api.getTrack()
	if err != nil {
		return "", err
	}

	down, err := createDownload(track)
	if err != nil {
		slog.Info("Error creating download: " + err.Error())
		return "", err
	}

	return down.Process()
}
