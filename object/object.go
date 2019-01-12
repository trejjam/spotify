package object

import (
	"fmt"
	"strings"
	"time"
)

type Device struct {
	Id               string
	IsActive         bool `json:"is_active"`
	IsPrivateSession bool `json:"is_private_session"`
	IsRestricted     bool `json:"is_restricted"`
	Name             string
	Type             string
	VolumePercent    int `json:"volume_percent"`
}

func (device *Device) String() string {
	return fmt.Sprintf(`{
  Id: %s
  IsActive: %t
  IsPrivateSession: %t
  IsRestricted: %t
  Name: %s
  Type: %s
  VolumePercent: %d
}`,
		device.Id,
		device.IsActive,
		device.IsPrivateSession,
		device.IsRestricted,
		device.Name,
		device.Type,
		device.VolumePercent,
	)
}

type Devices struct {
	Devices []Device
}

func (devices *Devices) String() string {
	var str strings.Builder
	str.WriteString("[\n")
	for i := 0; i < len(devices.Devices); i++ {
		device := devices.Devices[i]

		str.WriteString(device.String())
		str.WriteString(",\n")
	}
	str.WriteString("]\n")

	return str.String()
}

type TransferPlayback struct {
	Play    bool     `json:"play"`
	Devices []string `json:"device_ids"`
}

type PlayItems struct {
	Items   []PlayHistory
	Next    string
	Cursors Cursors
	Limit   int
	Href    string
}

type Cursors struct {
	After  string
	Before string
}

func (playHistories PlayItems) String() string {
	var str strings.Builder
	str.WriteString("[\n")
	for i := 0; i < len(playHistories.Items); i++ {
		item := playHistories.Items[i]

		str.WriteString(item.String())
		str.WriteString(",\n")
	}
	str.WriteString("]\n")

	return str.String()
}

type PlayHistory struct {
	Track    SimplifiedTrack
	PlayedAt time.Time `json:"played_at, string"`
	Context  Context
}

func (playHistory PlayHistory) String() string {
	return fmt.Sprintf(`{
  Track: %s
  PlayedAt: %s
  Context: %s
}`,
		playHistory.Track.String(),
		playHistory.PlayedAt,
		playHistory.Context.String(),
	)
}

type SimplifiedTrack struct {
	Artists          []SimplifiedArtist
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber       int      `json:"disc_number"`
	DurationInMs     int      `json:"duration_ms"`
	Explicit         bool
	ExternalUrls     ExternalUrl `json:"external_urls"`
	Href             string
	Id               string
	IsPlayable       bool         `json:"is_playable"`
	LinkedFrom       *LinkedTrack `json:"linked_from"`
	Name             string
	PreviewUrl       string `json:"preview_url"`
	TrackNumber      int    `json:"track_number"`
	Type             string
	Uri              string
}

func (simplifiedTrack SimplifiedTrack) String() string {
	var artistStr strings.Builder
	artistStr.WriteString("[\n")
	for i := 0; i < len(simplifiedTrack.Artists); i++ {
		artist := simplifiedTrack.Artists[i]

		artistStr.WriteString(artist.String())
		artistStr.WriteString(",\n")
	}
	artistStr.WriteString("]")

	linkedTrack := simplifiedTrack.LinkedFrom

	var linkedTrackString string
	if linkedTrack != nil {
		linkedTrackString = linkedTrack.String()
	}

	return fmt.Sprintf(`{
  Artists: %s
  AvailableMarkets: %s
  DiscNumber: %d
  DurationInMs: %d
  Explicit: %t
  ExternalUrls: %s
  Href: %s
  Id: %s
  IsPlayable: %t
  LinkedFrom: %s
  Name: %s
  PreviewUrl: %s
  TrackNumber: %d
  Type: %s
  Uri: %s
}`,
		artistStr.String(),
		simplifiedTrack.AvailableMarkets,
		simplifiedTrack.DiscNumber,
		simplifiedTrack.DurationInMs,
		simplifiedTrack.Explicit,
		simplifiedTrack.ExternalUrls,
		simplifiedTrack.Href,
		simplifiedTrack.Id,
		simplifiedTrack.IsPlayable,
		linkedTrackString,
		simplifiedTrack.Name,
		simplifiedTrack.PreviewUrl,
		simplifiedTrack.TrackNumber,
		simplifiedTrack.Type,
		simplifiedTrack.Uri,
	)
}

type SimplifiedArtist struct {
	ExternalUrls ExternalUrl `json:"external_urls"`
	Href         string
	Id           string
	Name         string
	Type         string
	Uri          string
}

func (simplifiedArtist SimplifiedArtist) String() string {
	return fmt.Sprintf(`{
  ExternalUrls: %s
  Href: %s
  Id: %s
  Name: %s
  Type: %s
  Uri: %s
}`,
		simplifiedArtist.ExternalUrls.String(),
		simplifiedArtist.Href,
		simplifiedArtist.Id,
		simplifiedArtist.Name,
		simplifiedArtist.Type,
		simplifiedArtist.Uri,
	)
}

type ExternalUrl struct {
	Spotify string
}

func (externalUrl ExternalUrl) String() string {
	return fmt.Sprintf("{Spotify: %s}",
		externalUrl.Spotify,
	)
}

type LinkedTrack struct {
	ExternalUrls ExternalUrl `json:"external_urls"`
	Href         string
	Id           string
	Type         string
	Uri          string
}

func (linkedTrack LinkedTrack) String() string {
	return fmt.Sprintf(`{
  ExternalUrls: %s
  Href: %s
  Id: %s
  Type: %s
  Uri: %s
}`,
		linkedTrack.ExternalUrls.String(),
		linkedTrack.Href,
		linkedTrack.Id,
		linkedTrack.Type,
		linkedTrack.Uri,
	)
}

type Context struct {
	Uri          string
	ExternalUrls ExternalUrl `json:"external_urls"`
	Href         string
	Type         string
}

func (context Context) String() string {
	return fmt.Sprintf(`{
  Uri: %s
  ExternalUrls: %s
  Href: %s
  Type: %s
}`,
		context.Uri,
		context.ExternalUrls.String(),
		context.Href,
		context.Type,
	)
}

type RepeatMode string

type CurrentlyPlaying struct {
	Timestamp            int
	Device               Device
	ProgresInMs          int    `json:"progress_ms"`
	IsPlaying            bool   `json:"is_playing"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Item                 Track
	ShuffleState         bool       `json:"shuffle_state"`
	RepeatState          RepeatMode `json:"repeat_state"`
	Context              Context
}

type CurrentlyPlayingTrack struct {
	Timestamp            int
	Context              Context
	ProgresInMs          int `json:"progress_ms"`
	Item                 Track
	CurrentlyPlayingType string `json:"currently_playing_type"`
	IsPlaying            bool   `json:"is_playing"`
}

type Track struct {
	Album            SimplifiedAlbum
	Artists          []Artist
	AvailableMarkets *[]string `json:"available_markets"`
	DiscNumber       int       `json:"disc_number"`
	DurationInMs     int       `json:"duration_ms"`
	Explicit         bool
	ExternalIds      ExternalId  `json:"external_ids"`
	ExternalUrls     ExternalUrl `json:"external_urls"`
	Href             string
	Id               string
	IsPlayable       bool   `json:"is_playable"`
	LinkedFrom       *Track `json:"linked_from"`
	Restrictions     *[]TrackRestriction
	Name             string
	Popularity       int
	PreviewUrl       string `json:"preview_url"`
	track_number     int
	Type             string
	Uri              string
}

type SimplifiedAlbum struct {
	AlbumGroup       string `json:"album_group"`
	AlbumType        string `json:"album_type"`
	Artists          []SimplifiedArtist
	AvailableMarkets []string    `json:"available_markets"`
	ExternalUrls     ExternalUrl `json:"external_urls"`
	Href             string
	Id               string
	Images           []Image
	Name             string
	Type             string
	Uri              string
}

type Artist struct {
	External_urls ExternalUrl `json:"external_urls"`
	Followers     *Followers
	Genres        []string
	Href          string
	Id            string
	Images        []Image
	Name          string
	Popularity    int
	Type          string
	Uri           string
}

type ExternalId struct {
	Isrc string
	Ean  string
	Upc  string
}

type TrackRestriction struct {
	Reason string
}

type Followers struct {
	//TODO not specified in public doc
}

type Image struct {
	Height int
	Url    string
	Width  int
}
