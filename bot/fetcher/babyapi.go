/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package dl

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"
	"github.com/BabiesIQ/SPOTIFY_MUSIC/config"
	babiesiq "github.com/BabiesIQ/_metaAPI/sdk/go"
)

var (
	babyAPISongSearch = func(ctx context.Context, query string) (*babiesiq.Song, error) {
		client, err := newBabiesIQClient()
		if err != nil {
			return nil, err
		}
		return client.Songs.Search(ctx, query, nil)
	}
	babyAPIVideoSearch = func(ctx context.Context, query string) (*babiesiq.Video, error) {
		client, err := newBabiesIQClient()
		if err != nil {
			return nil, err
		}
		return client.Videos.Search(ctx, query, nil)
	}
	babyAPISongDownload = func(ctx context.Context, query, destination string) (*babiesiq.SongDownloadResult, error) {
		client, err := newBabiesIQClient()
		if err != nil {
			return nil, err
		}
		return client.Songs.Download(ctx, query, destination, nil)
	}
	babyAPIVideoDownload = func(ctx context.Context, query, destination string) (*babiesiq.VideoDownloadResult, error) {
		client, err := newBabiesIQClient()
		if err != nil {
			return nil, err
		}
		return client.Videos.Download(ctx, query, destination, nil)
	}
)

func newBabiesIQClient() (*babiesiq.Client, error) {
	if strings.TrimSpace(config.ApiKey) == "" {
		return nil, fmt.Errorf("BabyAPI API_KEY is not configured")
	}

	return babiesiq.New(config.ApiKey, babiesiq.Config{
		BaseURL: strings.TrimRight(config.ApiUrl, "/"),
		Timeout: 30 * time.Second,
	})
}

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

func (y *youTubeData) searchWithBabyAPI(video bool) (utils.PlatformTracks, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// Resolve YouTube URLs to bare video IDs before handing off to the SDK.
	//
	// Root cause: when the SDK receives a full YouTube URL (e.g.
	// "https://youtu.be/DS-raAyMxl4"), it calls the YouTube InnerTube player
	// API (youtubei/v1/player) internally to resolve the URL before contacting
	// BabyAPI. Heroku datacenter IPs are blocked by YouTube's bot-detection on
	// that endpoint, so the SDK never reaches BabyAPI and returns "video not
	// found or unavailable".
	//
	// Fix: extract the bare 11-character video ID from the URL and pass that
	// directly. The SDK can then query BabyAPI without needing a YouTube
	// pre-flight, and the fallback to yt-dlp (which uses cookies) is preserved
	// if BabyAPI itself fails.
	query := y.Query
	if videoID := parseVideoID(normalizeYouTubeURL(y.Query)); videoID != "" {
		query = videoID
	}

	if video {
		result, err := babyAPIVideoSearch(ctx, query)
		if err != nil {
			return utils.PlatformTracks{}, fmt.Errorf("BabyAPI video search failed: %w", err)
		}
		if result == nil || result.VideoID == "" {
			return utils.PlatformTracks{}, fmt.Errorf("BabyAPI returned no video for %q", query)
		}
		return utils.PlatformTracks{Results: []utils.MusicTrack{babyAPIVideoTrack(result)}}, nil
	}

	result, err := babyAPISongSearch(ctx, query)
	if err != nil {
		return utils.PlatformTracks{}, fmt.Errorf("BabyAPI song search failed: %w", err)
	}
	if result == nil || result.VideoID == "" {
		return utils.PlatformTracks{}, fmt.Errorf("BabyAPI returned no song for %q", query)
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
