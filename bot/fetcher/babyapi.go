/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */
// Why the SDK is bypassed for BabiesIQ API calls
// ─────────────────────────────────────────────────────────────────────────────
// The official Go SDK (github.com/BabiesIQ/_metaAPI/sdk/go) calls
// resolveVideoID before every Songs/Videos request. resolveVideoID routes ANY
// 11-character string — including bare video IDs returned by our own YouTube
// search — through ytGetVideo → youtubei/v1/player, which is blocked on Heroku
// datacenter IPs by YouTube's bot-detection. This means the SDK is completely
// unusable on Heroku regardless of what query format we supply.
//
// Fix: call the BabiesIQ HTTP API directly, bypassing the SDK's resolveVideoID:
//   • URL / bare video ID queries: metadata via YouTube search API
//     (youtubei/v1/search — NOT blocked on Heroku), then POST to BabiesIQ.
//   • Text queries: same YouTube search path → video ID → BabiesIQ.
//
// The SDK package is still imported because our public types (Song, Video,
// SongDownloadResult, VideoDownloadResult) originate there and are used by
// tests and the rest of the bot.

package dl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/config"
	babiesiq "github.com/BabiesIQ/_metaAPI/sdk/go"
)

// ─── BabyAPI var-functions (replaceable in tests) ────────────────────────────

var (
	babyAPISongSearch = func(ctx context.Context, query string) (*babiesiq.Song, error) {
		return directSongLookup(ctx, query)
	}
	babyAPIVideoSearch = func(ctx context.Context, query string) (*babiesiq.Video, error) {
		return directVideoLookup(ctx, query)
	}
	babyAPISongDownload = func(ctx context.Context, query, destination string) (*babiesiq.SongDownloadResult, error) {
		return directSongDownload(ctx, query, destination)
	}
	babyAPIVideoDownload = func(ctx context.Context, query, destination string) (*babiesiq.VideoDownloadResult, error) {
		return directVideoDownload(ctx, query, destination)
	}
)

// ─── BabiesIQ direct HTTP client ─────────────────────────────────────────────

// biqAPIResp is the flat JSON envelope returned by the BabiesIQ HTTP API.
type biqAPIResp struct {
	Query    string `json:"query"`
	Status   string `json:"status"`
	Stream   string `json:"stream"`
	StreamID string `json:"stream_id"`
	Type     string `json:"type"`
	Error    string `json:"error"`
}

// biqRequest calls a BabiesIQ API endpoint directly via HTTP GET, completely
// bypassing the SDK and its internal YouTube player API pre-flight.
func biqRequest(ctx context.Context, endpoint, videoID string, download bool) (*biqAPIResp, error) {
	if strings.TrimSpace(config.ApiKey) == "" {
		return nil, fmt.Errorf("BabyAPI API_KEY is not configured")
	}

	baseURL := strings.TrimRight(config.ApiUrl, "/")
	params := url.Values{
		"query": {videoID},
		"api":   {config.ApiKey},
	}
	if download {
		params.Set("download", "true")
	}

	fullURL := baseURL + endpoint + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build BabiesIQ request: %w", err)
	}
	req.Header.Set("User-Agent", "biq-api-go/2.1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BabiesIQ request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to read BabiesIQ response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, fmt.Errorf("invalid or missing BabiesIQ API key")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("BabiesIQ rate limit exceeded")
	case http.StatusBadRequest:
		var errEnv struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(body, &errEnv) == nil && errEnv.Error != "" {
			return nil, fmt.Errorf("BabiesIQ API error: %s", errEnv.Error)
		}
		return nil, fmt.Errorf("BabiesIQ bad request: %s", string(body))
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("BabiesIQ server error %d: %s", resp.StatusCode, string(body))
	}

	var result biqAPIResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode BabiesIQ response: %w", err)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("BabiesIQ API error: %s", result.Error)
	}
	if result.Stream == "" {
		return nil, fmt.Errorf("BabiesIQ returned no stream URL for: %s", videoID)
	}
	return &result, nil
}

// ─── CDN polling + download ───────────────────────────────────────────────────

// cdnPollAndDownload polls streamURL every 2 s until the CDN signals the file
// is ready, then downloads it to destPath. Mirrors the SDK's pollAndDownload.
func cdnPollAndDownload(ctx context.Context, streamURL, destPath string, pollTimeout time.Duration) error {
	pollClient := &http.Client{Timeout: 15 * time.Second}
	deadline := time.Now().Add(pollTimeout)

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("CDN stream not ready after %s: %s", pollTimeout, streamURL)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "biq-api-go/2.1.0")

		resp, err := pollClient.Do(req)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
			}
			continue
		}
		status := resp.StatusCode
		_ = resp.Body.Close()

		switch {
		case status == http.StatusOK || status == http.StatusPartialContent:
			return cdnDownloadFile(ctx, streamURL, destPath)

		case status == http.StatusNoContent || status == 423 ||
			status == http.StatusNotFound || status == http.StatusGone:
			// Not ready yet — keep polling.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
			}

		case status == http.StatusUnauthorized || status == http.StatusForbidden:
			return fmt.Errorf("CDN access blocked (HTTP %d) — check API key", status)

		default:
			return fmt.Errorf("unexpected CDN status %d for: %s", status, streamURL)
		}
	}
}

// cdnDownloadFile streams the CDN content to destPath.
func cdnDownloadFile(ctx context.Context, streamURL, destPath string) error {
	dlClient := &http.Client{Timeout: 10 * time.Minute}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "biq-api-go/2.1.0")

	resp, err := dlClient.Do(req)
	if err != nil {
		return fmt.Errorf("CDN download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected CDN download status: %d", resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", destPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = os.Remove(destPath)
		return fmt.Errorf("CDN download interrupted: %w", err)
	}
	return nil
}

// ─── Video-ID resolution (no player API) ─────────────────────────────────────

// isBareVideoID reports whether s is a bare 11-character YouTube video ID
// (letters, digits, _ and - only). The bot stores bare IDs in the queue URL
// field, so we must recognise them without a URL prefix.
func isBareVideoID(s string) bool {
	if len(s) != 11 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// resolveVideoMeta resolves a query to a YouTube video ID and its metadata
// using only the YouTube search API (youtubei/v1/search). The player API
// (youtubei/v1/player) is intentionally never called — it is blocked on Heroku.
//
// query may be:
//   - A bare 11-char YouTube video ID (e.g. "dQw4w9WgXcQ" stored in cache URL)
//   - Any YouTube URL (youtu.be, watch?v=, /shorts/, …)
//   - A plain text search string
func resolveVideoMeta(ctx context.Context, query string) (videoID, title, artist, thumbnail string, durationSec int, err error) {
	// Case 1: URL or bare video ID — extract ID, fetch metadata via search API.
	vid := parseVideoID(normalizeYouTubeURL(query))
	if vid == "" && isBareVideoID(query) {
		vid = query // bare ID stored in the queue cache URL field
	}

	if vid != "" {
		videoID = vid
		// Try YouTube search with the video ID to get full metadata.
		if tracks, serr := youtubeSearch(vid, 10); serr == nil {
			for _, t := range tracks {
				if t.Id == vid {
					return vid, t.Title, t.Channel, t.Thumbnail, t.Duration, nil
				}
			}
		}
		// Fallback: oEmbed gives us the title, then search by title for full metadata.
		if oembedTitle, oerr := fetchYouTubeTitleFromOEmbed(vid); oerr == nil && oembedTitle != "" {
			if tracks, serr := youtubeSearch(oembedTitle, 10); serr == nil {
				for _, t := range tracks {
					if t.Id == vid {
						return vid, t.Title, t.Channel, t.Thumbnail, t.Duration, nil
					}
				}
			}
			return vid, oembedTitle, "", "", 0, nil
		}
		// Last resort: return the video ID as the title with no other metadata.
		return vid, vid, "", "", 0, nil
	}

	// Case 2: plain text query — use YouTube search API to find the best match.
	tracks, serr := youtubeSearch(query, 5)
	if serr != nil {
		return "", "", "", "", 0, fmt.Errorf("YouTube search failed for %q: %w", query, serr)
	}
	if len(tracks) == 0 {
		return "", "", "", "", 0, fmt.Errorf("no YouTube results for %q", query)
	}
	t := tracks[0]
	return t.Id, t.Title, t.Channel, t.Thumbnail, t.Duration, nil
}

// ─── Direct BabiesIQ search ───────────────────────────────────────────────────

func directSongLookup(ctx context.Context, query string) (*babiesiq.Song, error) {
	videoID, title, artist, thumbnail, durationSec, err := resolveVideoMeta(ctx, query)
	if err != nil {
		return nil, err
	}

	biqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	apiResp, err := biqRequest(biqCtx, "/api/song", videoID, false)
	if err != nil {
		return nil, err
	}

	return &babiesiq.Song{
		VideoID:   videoID,
		Title:     title,
		Artist:    artist,
		Thumbnail: thumbnail,
		Duration:  durationSec,
		StreamURL: apiResp.Stream,
	}, nil
}

func directVideoLookup(ctx context.Context, query string) (*babiesiq.Video, error) {
	videoID, title, artist, thumbnail, durationSec, err := resolveVideoMeta(ctx, query)
	if err != nil {
		return nil, err
	}

	biqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	apiResp, err := biqRequest(biqCtx, "/api/video", videoID, false)
	if err != nil {
		return nil, err
	}

	return &babiesiq.Video{
		VideoID:   videoID,
		Title:     title,
		Channel:   artist,
		Thumbnail: thumbnail,
		Duration:  durationSec,
		StreamURL: apiResp.Stream,
	}, nil
}

// ─── Direct BabiesIQ download ─────────────────────────────────────────────────

func directSongDownload(ctx context.Context, videoID, destPath string) (*babiesiq.SongDownloadResult, error) {
	if strings.TrimSpace(videoID) == "" {
		return nil, fmt.Errorf("BabiesIQ song download requires a video ID")
	}

	biqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	apiResp, err := biqRequest(biqCtx, "/api/song", videoID, true)
	if err != nil {
		return nil, err
	}

	if err := cdnPollAndDownload(ctx, apiResp.Stream, destPath, 2*time.Minute); err != nil {
		return nil, err
	}

	return &babiesiq.SongDownloadResult{
		Song:     &babiesiq.Song{VideoID: videoID, StreamURL: apiResp.Stream},
		FilePath: destPath,
	}, nil
}

func directVideoDownload(ctx context.Context, videoID, destPath string) (*babiesiq.VideoDownloadResult, error) {
	if strings.TrimSpace(videoID) == "" {
		return nil, fmt.Errorf("BabiesIQ video download requires a video ID")
	}

	biqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	apiResp, err := biqRequest(biqCtx, "/api/video", videoID, true)
	if err != nil {
		return nil, err
	}

	if err := cdnPollAndDownload(ctx, apiResp.Stream, destPath, 3*time.Minute); err != nil {
		return nil, err
	}

	return &babiesiq.VideoDownloadResult{
		Video:    &babiesiq.Video{VideoID: videoID, StreamURL: apiResp.Stream},
		FilePath: destPath,
	}, nil
}

// ─── Track mapping helpers ────────────────────────────────────────────────────

func babyAPISongTrack(song *babiesiq.Song) utils.MusicTrack {
	return utils.MusicTrack{
		Title:     song.Title,
		Id:        song.VideoID,
		Url:       song.VideoID,
		Thumbnail: song.Thumbnail,
		Duration:  song.Duration,
		Channel:   song.Artist,
		Platform:  utils.YouTube,
	}
}

func babyAPIVideoTrack(video *babiesiq.Video) utils.MusicTrack {
	return utils.MusicTrack{
		Title:     video.Title,
		Id:        video.VideoID,
		Url:       video.VideoID,
		Thumbnail: video.Thumbnail,
		Duration:  video.Duration,
		Channel:   video.Channel,
		Platform:  utils.YouTube,
	}
}

func babyAPIQuery(query string, videoID string) string {
	if strings.TrimSpace(videoID) != "" {
		return videoID
	}
	return strings.TrimSpace(query)
}

// ─── youTubeData BabyAPI methods ─────────────────────────────────────────────

func (y *youTubeData) searchWithBabyAPI(video bool) (utils.PlatformTracks, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	if video {
		result, err := babyAPIVideoSearch(ctx, y.Query)
		if err != nil {
			return utils.PlatformTracks{}, fmt.Errorf("BabyAPI video search failed: %w", err)
		}
		if result == nil || result.VideoID == "" {
			return utils.PlatformTracks{}, fmt.Errorf("BabyAPI returned no video for %q", y.Query)
		}
		return utils.PlatformTracks{Results: []utils.MusicTrack{babyAPIVideoTrack(result)}}, nil
	}

	result, err := babyAPISongSearch(ctx, y.Query)
	if err != nil {
		return utils.PlatformTracks{}, fmt.Errorf("BabyAPI song search failed: %w", err)
	}
	if result == nil || result.VideoID == "" {
		return utils.PlatformTracks{}, fmt.Errorf("BabyAPI returned no song for %q", y.Query)
	}
	return utils.PlatformTracks{Results: []utils.MusicTrack{babyAPISongTrack(result)}}, nil
}

func (y *youTubeData) getInfoWithBabyAPI(video bool) (utils.PlatformTracks, error) {
	if strings.TrimSpace(y.Query) == "" {
		return utils.PlatformTracks{}, fmt.Errorf("the query is empty")
	}
	return y.searchWithBabyAPI(video)
}

func (y *youTubeData) getTrackWithBabyAPI(video bool) (utils.TrackInfo, error) {
	tracks, err := y.getInfoWithBabyAPI(video)
	if err != nil {
		return utils.TrackInfo{}, err
	}
	if len(tracks.Results) == 0 {
		return utils.TrackInfo{}, fmt.Errorf("BabyAPI returned no tracks for %q", y.Query)
	}

	track := tracks.Results[0]
	return utils.TrackInfo{
		Id:       track.Id,
		URL:      track.Url,
		Platform: utils.YouTube,
	}, nil
}

func (y *youTubeData) downloadWithBabyAPI(videoID string, video bool) (string, error) {
	if strings.TrimSpace(videoID) == "" {
		return "", fmt.Errorf("BabyAPI download requires a YouTube video ID")
	}

	extension := "mp3"
	if video {
		extension = "mp4"
	}
	destination := filepath.Join(config.DownloadsDir, fmt.Sprintf("%s.%s", videoID, extension))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	query := babyAPIQuery(y.Query, videoID)

	if video {
		result, err := babyAPIVideoDownload(ctx, query, destination)
		if err != nil {
			return "", fmt.Errorf("BabyAPI video download failed: %w", err)
		}
		if result == nil || result.FilePath == "" {
			return "", fmt.Errorf("BabyAPI video download returned no file")
		}
		return result.FilePath, nil
	}

	result, err := babyAPISongDownload(ctx, query, destination)
	if err != nil {
		return "", fmt.Errorf("BabyAPI song download failed: %w", err)
	}
	if result == nil || result.FilePath == "" {
		return "", fmt.Errorf("BabyAPI song download returned no file")
	}
	return result.FilePath, nil
}
