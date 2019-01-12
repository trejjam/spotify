package main

import (
	"fmt"
	. "github.com/trejjam/spotify/accessToken"
	. "github.com/trejjam/spotify/player"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Errorf("Use username and password as arguments"))
		return
	}

	username := os.Args[1]
	password := os.Args[2]

	accessToken, err := GetAccessToken(username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println(accessToken.AccessToken)

	devices, err := GetDevices(accessToken)
	if err != nil {
		panic(err)
	}

	fmt.Println(devices.String())

	now := time.Now()

	recentlyPlayedTracks, err := RecentlyPlayedTracksAfter(accessToken, 1, now.Add(-time.Hour))
	if err != nil {
		panic(err)
	}

	fmt.Println(recentlyPlayedTracks.String())

	currentlyPlaying, err := GetCurrentlyPlaying(accessToken, "CZ")
	if err != nil {
		panic(err)
	}

	fmt.Println(currentlyPlaying)

	currentlyPlayingTrack, err := GetCurrentlyPlayingTrack(accessToken, "CZ")
	if err != nil {
		panic(err)
	}

	fmt.Println(currentlyPlayingTrack)
}
