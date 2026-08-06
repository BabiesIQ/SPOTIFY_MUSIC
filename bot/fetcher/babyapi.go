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
