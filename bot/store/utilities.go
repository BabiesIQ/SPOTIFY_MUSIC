/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package db

import (
	"fmt"
)

// convertToKey converts an int64 ID into a string format suitable for use as a cache key.
func convertToKey(id int64) string {
	return fmt.Sprintf("%d", id)
}

// sliceContains checks if a given int64 slice sliceContains a specific ID.
// It returns true if the ID is found, and false otherwise.
func sliceContains(list []int64, id int64) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}
	return false
}
