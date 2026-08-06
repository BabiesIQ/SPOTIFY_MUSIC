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
	"errors"
	"testing"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/utils"
	babiesiq "github.com/BabiesIQ/_metaAPI/sdk/go"
)

func withBabyAPIMocks(t *testing.T) {
	t.Helper()
	originalSongSearch := babyAPISongSearch
	originalVideoSearch := babyAPIVideoSearch
	originalSongDownload := babyAPISongDownload
	originalVideoDownload := babyAPIVideoDownload
	originalYouTubeSearch := youtubeSearch
	originalYtDlpDownload := ytDlpDownload
	t.Cleanup(func() {
		babyAPISongSearch = originalSongSearch
		babyAPIVideoSearch = originalVideoSearch
		babyAPISongDownload = originalSongDownload
		babyAPIVideoDownload = originalVideoDownload
		youtubeSearch = originalYouTubeSearch
		ytDlpDownload = originalYtDlpDownload
	})
}

func TestBabyAPISongTrackMapping(t *testing.T) {
	track := babyAPISongTrack(&babiesiq.Song{
		Title:     "Test song",
		VideoID:   "dQw4w9WgXcQ",
		Thumbnail: "https://img.example/song.jpg",
		Duration:  213,
		Artist:    "Test artist",
	})

	if track.Title != "Test song" ||
		track.Id != "dQw4w9WgXcQ" ||
		track.Url != "dQw4w9WgXcQ" ||
		track.Thumbnail != "https://img.example/song.jpg" ||
		track.Duration != 213 ||
		track.Channel != "Test artist" ||
		track.Platform != utils.YouTube {
		t.Fatalf("unexpected song mapping: %+v", track)
	}
}

func TestBabyAPIVideoTrackMapping(t *testing.T) {
	track := babyAPIVideoTrack(&babiesiq.Video{
		Title:     "Test video",
		VideoID:   "dQw4w9WgXcQ",
		Thumbnail: "https://img.example/video.jpg",
		Duration:  240,
		Channel:   "Test channel",
	})

	if track.Title != "Test video" ||
		track.Id != "dQw4w9WgXcQ" ||
		track.Url != "dQw4w9WgXcQ" ||
		track.Thumbnail != "https://img.example/video.jpg" ||
		track.Duration != 240 ||
		track.Channel != "Test channel" ||
		track.Platform != utils.YouTube {
		t.Fatalf("unexpected video mapping: %+v", track)
	}
}

func TestYouTubeQueryValidationAcceptsNamesAndRejectsOtherURLs(t *testing.T) {
	if !createYouTubeData("Tu Hi Rab Tu Hi Dua").isValid() {
		t.Fatal("plain song names should be handled by the YouTube service")
	}
	if !createYouTubeData("https://www.youtube.com/watch?v=dQw4w9WgXcQ").isValid() {
		t.Fatal("YouTube URLs should be handled by the YouTube service")
	}
	if createYouTubeData("https://open.spotify.com/track/example").isValid() {
		t.Fatal("Spotify URLs should not be handled by the YouTube service")
	}
}

// TestBabyAPISearchReceivesRawQuery verifies that searchWithBabyAPI passes the
// raw query (URL or text) directly to the babyAPISongSearch / babyAPIVideoSearch
// var functions. URL-to-video-ID resolution is the responsibility of the var
// function itself (directSongLookup / directVideoLookup), so the mock captures
// whatever the caller supplies.
func TestBabyAPISearchReceivesRawQuery(t *testing.T) {
	withBabyAPIMocks(t)
	var songQuery, videoQuery string
	babyAPISongSearch = func(_ context.Context, query string) (*babiesiq.Song, error) {
		songQuery = query
		return &babiesiq.Song{Title: "song", VideoID: "song-id"}, nil
	}
	babyAPIVideoSearch = func(_ context.Context, query string) (*babiesiq.Video, error) {
		videoQuery = query
		return &babiesiq.Video{Title: "video", VideoID: "video-id"}, nil
	}

	y := createYouTubeData("https://youtu.be/dQw4w9WgXcQ")
	audio, err := y.searchForPlayback(false)
	if err != nil {
		t.Fatal(err)
	}
	video, err := y.searchForPlayback(true)
	if err != nil {
		t.Fatal(err)
	}

	// The raw query (full URL) must be forwarded unchanged — resolution is
	// done inside directSongLookup / directVideoLookup, not by the caller.
	if songQuery != y.Query || videoQuery != y.Query {
		t.Fatalf("var functions must receive raw query %q, got: song=%q video=%q",
			y.Query, songQuery, videoQuery)
	}
	if audio.Results[0].Id != "song-id" || video.Results[0].Id != "video-id" {
		t.Fatalf("unexpected results: audio=%+v video=%+v", audio.Results, video.Results)
	}
}

func TestBabyAPIDownloadUsesAudioAndVideoEndpoints(t *testing.T) {
	withBabyAPIMocks(t)
	var songQuery, songPath, videoQuery, videoPath string
	babyAPISongDownload = func(_ context.Context, query, destination string) (*babiesiq.SongDownloadResult, error) {
		songQuery, songPath = query, destination
		return &babiesiq.SongDownloadResult{Song: &babiesiq.Song{VideoID: query}, FilePath: destination}, nil
	}
	babyAPIVideoDownload = func(_ context.Context, query, destination string) (*babiesiq.VideoDownloadResult, error) {
		videoQuery, videoPath = query, destination
		return &babiesiq.VideoDownloadResult{Video: &babiesiq.Video{VideoID: query}, FilePath: destination}, nil
	}

	y := createYouTubeData("song name")
	audio, err := y.downloadTrack(utils.TrackInfo{Id: "audio-id"}, false)
	if err != nil {
		t.Fatal(err)
	}
	video, err := y.downloadTrack(utils.TrackInfo{Id: "video-id"}, true)
	if err != nil {
		t.Fatal(err)
	}

	if songQuery != "audio-id" || videoQuery != "video-id" ||
		songPath != audio || videoPath != video ||
		!endsWithSuffix(audio, "audio-id.mp3") ||
		!endsWithSuffix(video, "video-id.mp4") {
		t.Fatalf("unexpected download calls: song=(%q,%q) video=(%q,%q)", songQuery, songPath, videoQuery, videoPath)
	}
}

func endsWithSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func TestBabyAPISearchFailureFallsBackToYouTube(t *testing.T) {
	withBabyAPIMocks(t)
	babyAPISongSearch = func(context.Context, string) (*babiesiq.Song, error) {
		return nil, errors.New("BabyAPI unavailable")
	}
	youtubeSearch = func(_ string, _ int) ([]utils.MusicTrack, error) {
		return []utils.MusicTrack{{Id: "fallback-id", Title: "fallback"}}, nil
	}

	result, err := createYouTubeData("fallback query").searchForPlayback(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 || result.Results[0].Id != "fallback-id" {
		t.Fatalf("unexpected fallback result: %+v", result.Results)
	}
}

func TestBabyAPIVideoSearchFailureDoesNotUseAudioSearch(t *testing.T) {
	withBabyAPIMocks(t)
	babyAPIVideoSearch = func(context.Context, string) (*babiesiq.Video, error) {
		return nil, errors.New("video endpoint unavailable")
	}
	babyAPISongSearch = func(context.Context, string) (*babiesiq.Song, error) {
		t.Fatal("audio search must not be used as the video fallback")
		return nil, nil
	}
	youtubeSearch = func(_ string, _ int) ([]utils.MusicTrack, error) {
		return []utils.MusicTrack{{Id: "video-fallback-id", Title: "fallback"}}, nil
	}

	result, err := createYouTubeData("video query").searchForPlayback(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 || result.Results[0].Id != "video-fallback-id" {
		t.Fatalf("unexpected video fallback result: %+v", result.Results)
	}
}

func TestBabyAPIDownloadFailureFallsBackToYtDlp(t *testing.T) {
	withBabyAPIMocks(t)
	babyAPISongDownload = func(context.Context, string, string) (*babiesiq.SongDownloadResult, error) {
		return nil, errors.New("BabyAPI unavailable")
	}
	ytDlpDownload = func(_ *youTubeData, videoID string, video bool) (string, error) {
		if video {
			return videoID + ".mp4", nil
		}
		return videoID + ".mp3", nil
	}

	result, err := createYouTubeData("fallback query").downloadTrack(utils.TrackInfo{Id: "fallback-id"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if result != "fallback-id.mp3" {
		t.Fatalf("unexpected yt-dlp fallback result: %q", result)
	}
}

// TestIsBareVideoID verifies the bare video ID detector used by resolveVideoMeta
// to handle cached queue URLs that store a bare ID rather than a full URL.
func TestIsBareVideoID(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"dQw4w9WgXcQ", true},
		{"DS-raAyMxl4", true},
		{"TAHRNfVci48", true},
		{"ERfowTjYY-A", true},
		{"ZM8rAsTT7yE", true},
		{"https://youtu.be/dQw4w9WgXcQ", false}, // URL, not bare
		{"Shape of You", false},                  // text query
		{"dQw4w9WgXcQX", false},                  // 12 chars
		{"dQw4w9WgXc", false},                    // 10 chars
		{"dQw4w9WgXc!", false},                   // invalid char
	}
	for _, tc := range cases {
		if got := isBareVideoID(tc.input); got != tc.want {
			t.Errorf("isBareVideoID(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}
