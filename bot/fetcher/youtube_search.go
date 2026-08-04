/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package dl

import (
	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	ytBaseURL   = "https://www.youtube.com"
	ytWatchURL  = ytBaseURL + "/watch?v="
	ytAPIKey    = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
	ytClientVer = "2.20240229.01.00"
)

var (
	labelDurationRe = regexp.MustCompile(`(\d+)\s*(hours?|minutes?|seconds?)`)
	videoIDRe1      = regexp.MustCompile(`(?i)(?:youtube\.com/(?:watch\?v=|embed/|shorts/|live/)|youtu\.be/)([A-Za-z0-9_-]{11})`)
	videoIDRe2      = regexp.MustCompile(`(?:v=|\/)([0-9A-Za-z_-]{11})`)
	playlistIDRe1   = regexp.MustCompile(`(?i)(?:youtube\.com|music\.youtube\.com).*(?:\?|&)list=([A-Za-z0-9_-]+)`)
	playlistIDRe2   = regexp.MustCompile(`list=([0-9A-Za-z_-]+)`)
)

// youtubeContext returns the standard InnerTube context payload.
func youtubeContext() map[string]any {
	return map[string]any{
		"context": map[string]any{
			"client": map[string]any{
				"clientName":    "WEB",
				"clientVersion": ytClientVer,
			},
		},
	}
}

// youtubePost builds, sends, and decodes a POST request to a YouTube InnerTube endpoint.
// extraFields are merged into the top-level payload alongside "context".
func youtubePost(ctx context.Context, path string, extraFields map[string]any) (map[string]any, error) {
	payload := youtubeContext()
	for k, v := range extraFields {
		payload[k] = v
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	endpoint := ytBaseURL + path + "?key=" + ytAPIKey
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return nil, fmt.Errorf("youtube %s failed: status=%d body=%q", path, res.StatusCode, snippet)
	}

	var out map[string]any
	if err := json.NewDecoder(io.LimitReader(res.Body, 10*1024*1024)).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}

func searchYouTube(query string, limit int) ([]utils.MusicTrack, error) {
	payload := map[string]any{
		"context": map[string]any{
			"client": map[string]any{
				"clientName":    "WEB",
				"clientVersion": "2.20250101.01.00",
				"hl":            "en",
				"gl":            "IN",
			},
		},
		"query": query,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal search payload: %w", err)
	}

	endpoint := ytBaseURL + "/youtubei/v1/search?key=" + ytAPIKey
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("build search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("youtube search failed: status=%d %s body=%q",
			resp.StatusCode, resp.Status, raw)
	}

	var data map[string]any
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10*1024*1024)).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	root := deepLookup(data,
		"contents",
		"twoColumnSearchResultsRenderer",
		"primaryContents",
		"sectionListRenderer",
		"contents",
	)

	var tracks []utils.MusicTrack
	parseResults(root, &tracks, limit)
	return tracks, nil
}

func parseResults(node any, tracks *[]utils.MusicTrack, limit int) {
	stack := []any{node}
	for len(stack) > 0 && len(*tracks) < limit {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch v := curr.(type) {
		case []any:
			for i := len(v) - 1; i >= 0; i-- {
				stack = append(stack, v[i])
			}

		case map[string]any:
			if vr, ok := deepLookup(v, "videoRenderer").(map[string]any); ok {
				if isCurrentlyLive(vr) {
					continue
				}
				id := safeString(vr["videoId"])
				title := safeString(deepLookup(vr, "title", "runs", 0, "text"))
				durationText := safeString(deepLookup(vr, "lengthText", "simpleText"))
				if id == "" || title == "" || durationText == "" {
					continue
				}
				*tracks = append(*tracks, utils.MusicTrack{
					Id:        id,
					Url:       ytWatchURL + id,
					Title:     title,
					Thumbnail: safeString(deepLookup(vr, "thumbnail", "thumbnails", 0, "url")),
					Duration:  parseDuration(durationText),
					Views:     safeString(deepLookup(vr, "viewCountText", "simpleText")),
					Channel:   safeString(deepLookup(vr, "ownerText", "runs", 0, "text")),
					Platform:  utils.YouTube,
				})
				continue
			}

			for _, child := range v {
				stack = append(stack, child)
			}
		}
	}
}

// isCurrentlyLive reports whether a videoRenderer map carries the LIVE_NOW badge.
func isCurrentlyLive(vr map[string]any) bool {
	badges, ok := vr["badges"].([]any)
	if !ok {
		return false
	}
	for _, badge := range badges {
		meta, ok := deepLookup(badge, "metadataBadgeRenderer").(map[string]any)
		if !ok {
			continue
		}
		if safeString(meta["style"]) == "BADGE_STYLE_TYPE_LIVE_NOW" {
			return true
		}
	}
	return false
}

// fetchYouTubeTitleFromOEmbed fetches the video title using YouTube's oEmbed API.
func fetchYouTubeTitleFromOEmbed(videoID string) (string, error) {
	apiURL := fmt.Sprintf("https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=%s&format=json", videoID)

	resp, err := client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("oEmbed request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oEmbed returned status code: %d", resp.StatusCode)
	}

	var data struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1*1024*1024)).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode oEmbed response: %w", err)
	}

	if data.Title == "" {
		return "", errors.New("oEmbed response contained empty title")
	}

	return data.Title, nil
}

func fetchYouTubeVideo(ctx context.Context, videoID string) (utils.PlatformTracks, error) {
	resp, err := youtubePost(ctx, "/youtubei/v1/player", map[string]any{"videoId": videoID})
	if err != nil {
		return utils.PlatformTracks{}, err
	}

	video := convertPlayerToTrack(resp)
	if video.Id == "" {
		return utils.PlatformTracks{}, errors.New("video not found")
	}
	return utils.PlatformTracks{Results: []utils.MusicTrack{video}}, nil
}

func fetchYouTubePlaylist(ctx context.Context, playlistID string) (utils.PlatformTracks, error) {
	resp, err := youtubePost(ctx, "/youtubei/v1/browse", map[string]any{"browseId": "VL" + playlistID})
	if err != nil {
		return utils.PlatformTracks{}, err
	}

	videos := parsePlaylistVideos(resp)
	return constructTrackList(videos, convertYTVideo), nil
}

func FetchYouTubeMixPlaylist(ctx context.Context, playlistID string) (utils.PlatformTracks, error) {
	resp, err := youtubePost(ctx, "/youtubei/v1/next", map[string]any{"playlistId": playlistID})
	if err != nil {
		return utils.PlatformTracks{}, err
	}

	videos := parseMixPlaylistVideos(resp)
	return constructTrackList(videos, convertMixVideo), nil
}

// constructTrackList converts raw renderer maps to MusicTrack, dropping empty IDs.
func constructTrackList(videos []map[string]any, mapper func(map[string]any) utils.MusicTrack) utils.PlatformTracks {
	out := make([]utils.MusicTrack, 0, len(videos))
	for _, v := range videos {
		if t := mapper(v); t.Id != "" {
			out = append(out, t)
		}
	}
	return utils.PlatformTracks{Results: out}
}

func convertYTVideo(v map[string]any) utils.MusicTrack {
	id := deepLookupStr(v, "videoId")
	return utils.MusicTrack{
		Id:        id,
		Title:     deepLookupStr(v, "title", "runs", 0, "text"),
		Url:       ytWatchURL + id,
		Thumbnail: pickYTThumb(v),
		Channel:   deepLookupStr(v, "shortBylineText", "runs", 0, "text"),
		Duration:  parseYTDuration(v),
		Views:     deepLookupStr(v, "viewCountText", "simpleText"),
		Platform:  utils.YouTube,
	}
}

func convertMixVideo(v map[string]any) utils.MusicTrack {
	id := deepLookupStr(v, "videoId")
	return utils.MusicTrack{
		Id:        id,
		Title:     deepLookupStr(v, "title", "simpleText"),
		Url:       ytWatchURL + id,
		Thumbnail: pickYTThumb(v),
		Channel:   deepLookupStr(v, "shortBylineText", "runs", 0, "text"),
		Duration:  parseYTDuration(v),
		Platform:  utils.YouTube,
	}
}

func convertPlayerToTrack(src map[string]any) utils.MusicTrack {
	id := deepLookupStr(src, "videoDetails", "videoId")
	return utils.MusicTrack{
		Id:        id,
		Title:     deepLookupStr(src, "videoDetails", "title"),
		Url:       ytWatchURL + id,
		Thumbnail: pickYTPlayerThumb(src),
		Channel:   deepLookupStr(src, "videoDetails", "author"),
		Duration:  parseIntStr(deepLookupStr(src, "videoDetails", "lengthSeconds")),
		Views:     deepLookupStr(src, "videoDetails", "viewCount"),
		Platform:  utils.YouTube,
	}
}

func parsePlaylistVideos(src map[string]any) []map[string]any {
	contents := deepLookupArray(src,
		"contents",
		"twoColumnBrowseResultsRenderer",
		"tabs", 0,
		"tabRenderer",
		"content",
		"sectionListRenderer",
		"contents", 0,
		"itemSectionRenderer",
		"contents", 0,
		"playlistVideoListRenderer",
		"contents",
	)
	var out []map[string]any
	for _, c := range contents {
		if v, ok := c["playlistVideoRenderer"].(map[string]any); ok {
			out = append(out, v)
		}
	}
	return out
}

func parseMixPlaylistVideos(src map[string]any) []map[string]any {
	contents := deepLookupArray(src,
		"contents",
		"twoColumnWatchNextResults",
		"playlist", "playlist", "contents",
	)
	var out []map[string]any
	for _, c := range contents {
		if v, ok := c["playlistPanelVideoRenderer"].(map[string]any); ok {
			out = append(out, v)
		}
	}
	return out
}

func pickYTThumb(v map[string]any) string {
	return resolveLastThumbURL(deepLookupArray(v, "thumbnail", "thumbnails"))
}

func pickYTPlayerThumb(src map[string]any) string {
	return resolveLastThumbURL(deepLookupArray(src, "videoDetails", "thumbnail", "thumbnails"))
}

// resolveLastThumbURL returns the URL of the last (highest-res) thumbnail, or "".
func resolveLastThumbURL(thumbs []map[string]any) string {
	if len(thumbs) == 0 {
		return ""
	}
	t, _ := thumbs[len(thumbs)-1]["url"].(string)
	return t
}

func normalizeYouTubeURL(rawURL string) string {
	var id string
	switch {
	case strings.Contains(rawURL, "youtu.be/"):
		id = parseSegment(rawURL, "youtu.be/")
	case strings.Contains(rawURL, "youtube.com/shorts/"):
		id = parseSegment(rawURL, "youtube.com/shorts/")
	default:
		return rawURL
	}
	return ytWatchURL + id
}

// parseSegment splits on sep, then strips query string and fragment.
func parseSegment(u, sep string) string {
	after := strings.SplitN(u, sep, 2)[1]
	after = strings.SplitN(after, "?", 2)[0]
	after = strings.SplitN(after, "#", 2)[0]
	return after
}

func parseVideoID(u string) string {
	if m := videoIDRe1.FindStringSubmatch(u); len(m) > 1 {
		return m[1]
	}
	if m := videoIDRe2.FindStringSubmatch(u); len(m) > 1 {
		return m[1]
	}
	return ""
}

func parsePlaylistID(u string) string {
	if m := playlistIDRe1.FindStringSubmatch(u); len(m) > 1 {
		return m[1]
	}
	if m := playlistIDRe2.FindStringSubmatch(u); len(m) > 1 {
		return m[1]
	}
	return ""
}

func parseYTDuration(v map[string]any) int {
	if txt := deepLookupStr(v, "lengthText", "simpleText"); txt != "" {
		return parseTimeToSeconds(txt)
	}
	if label := deepLookupStr(v, "lengthText", "accessibility", "accessibilityData", "label"); label != "" {
		return parseLabelDuration(label)
	}
	return 0
}

// parseDuration handles "H:MM:SS" / "M:SS" / "SS" colon-separated strings.
func parseDuration(s string) int {
	parts := strings.Split(s, ":")
	total, mul := 0, 1
	for i := len(parts) - 1; i >= 0; i-- {
		total += parseIntStr(parts[i]) * mul
		mul *= 60
	}
	return total
}

// parseTimeToSeconds is like parseDuration but uses strconv and returns 0 on any error.
func parseTimeToSeconds(s string) int {
	parts := strings.Split(s, ":")
	total := 0
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return 0
		}
		total = total*60 + n
	}
	return total
}

func parseLabelDuration(s string) int {
	total := 0
	for _, m := range labelDurationRe.FindAllStringSubmatch(s, -1) {
		n, _ := strconv.Atoi(m[1])
		switch {
		case strings.HasPrefix(m[2], "hour"):
			total += n * 3600
		case strings.HasPrefix(m[2], "minute"):
			total += n * 60
		default:
			total += n
		}
	}
	return total
}

func deepLookup(v any, path ...any) any {
	cur := v
	for _, p := range path {
		if cur == nil {
			return nil
		}
		switch k := p.(type) {
		case string:
			m, ok := cur.(map[string]any)
			if !ok {
				return nil
			}
			cur = m[k]
		case int:
			a, ok := cur.([]any)
			if !ok || k < 0 || k >= len(a) {
				return nil
			}
			cur = a[k]
		}
	}
	return cur
}

func deepLookupStr(src any, path ...any) string {
	s, _ := deepLookup(src, path...).(string)
	return s
}

func deepLookupArray(src any, path ...any) []map[string]any {
	arr, ok := deepLookup(src, path...).([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(arr))
	for _, v := range arr {
		if m, ok := v.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func safeString(v any) string {
	s, _ := v.(string)
	return s
}

func parseIntStr(s string) int {
	n := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			n = n*10 + int(r-'0')
		}
	}
	return n
}
