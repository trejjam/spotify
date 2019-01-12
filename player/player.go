package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/trejjam/spotify/accessToken"
	"github.com/trejjam/spotify/object"
	"github.com/trejjam/spotify/spotifyError"
	"github.com/trejjam/spotify/url/player"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func initClient() *http.Client {
	return &http.Client{
		CheckRedirect: nil,
	}
}

var httpClient *http.Client

func getClient() *http.Client {
	if httpClient == nil {
		httpClient = initClient()
	}
	return httpClient
}

func performGet(accessToken *accessToken.AccessToken, url string, mapObject interface{}) error {
	request, err := accessToken.CreateGetRequest(url)
	if err != nil {
		return err
	}

	client := getClient()

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = spotifyError.Validate200Response(response)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(mapObject)
	if err != nil {
		return err
	}

	return nil
}

func performDebugGet(accessToken *accessToken.AccessToken, url string, mapObject interface{}) error {
	request, err := accessToken.CreateGetRequest(url)
	if err != nil {
		return err
	}

	client := getClient()

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = spotifyError.Validate200Response(response)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Print(string(body))

	err = json.Unmarshal(body, mapObject)
	if err != nil {
		return err
	}

	return nil
}

func performEmptyNoResponsePost(accessToken *accessToken.AccessToken, url string) error {
	request, err := accessToken.CreatePostRequest(url, nil)
	if err != nil {
		return err
	}

	client := getClient()

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = spotifyError.Validate204Response(response)
	if err != nil {
		return err
	}

	return nil
}

func performEmptyNoResponsePut(accessToken *accessToken.AccessToken, url string) error {
	return performNoResponsePut(accessToken, url, nil)
}

func performNoResponsePut(accessToken *accessToken.AccessToken, url string, bodyObject interface{}) error {
	var body io.Reader
	if bodyObject != nil {
		bodyObjectJson, err := json.Marshal(bodyObject)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(bodyObjectJson)
	}

	request, err := accessToken.CreatePutRequest(url, body)
	if err != nil {
		return err
	}

	client := getClient()

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if bodyObject == nil {
		err = spotifyError.Validate204Response(response)
		if err != nil {
			return err
		}
	} else {
		err = spotifyError.Validate2xxResponse(response)
		if err != nil {
			return err
		}
	}

	return nil
}

func NextTrack(accessToken *accessToken.AccessToken, device *object.Device) error {
	requestUrl, err := url.Parse(player.SkipUsersPlaybackToNextTrack)
	if err != nil {
		return err
	}

	if device != nil {
		requestQuery := requestUrl.Query()
		requestQuery.Add("device_id", device.Id)

		requestUrl.RawQuery = requestQuery.Encode()
	}

	err = performEmptyNoResponsePost(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func SeekTrack(accessToken *accessToken.AccessToken, positionInMs int, device *object.Device) error {
	requestUrl, err := url.Parse(player.SeekToPositionInCurrentlyPlayingTrack)
	if err != nil {
		return err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("position_ms", strconv.Itoa(positionInMs))
	if device != nil {
		requestQuery.Add("device_id", device.Id)
	}

	requestUrl.RawQuery = requestQuery.Encode()

	err = performEmptyNoResponsePut(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func GetDevices(accessToken *accessToken.AccessToken) (*object.Devices, error) {
	devices := new(object.Devices)
	err := performGet(accessToken, player.GetUserAvailableDevices, devices)

	if err != nil {
		return nil, err
	}

	return devices, nil
}

func ToggleShuffle(accessToken *accessToken.AccessToken, shuffle bool, device *object.Device) error {
	requestUrl, err := url.Parse(player.ToggleShuffleForUserPlayback)
	if err != nil {
		return err
	}

	shuffleString := "false"
	if shuffle {
		shuffleString = "true"
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("state", shuffleString)
	if device != nil {
		requestQuery.Add("device_id", device.Id)
	}

	requestUrl.RawQuery = requestQuery.Encode()

	err = performEmptyNoResponsePut(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func TransferPlayback(accessToken *accessToken.AccessToken, play bool, device *object.Device) error {
	if device == nil {
		return &spotifyError.RequiredParameterError{
			Method:    "ToggleShuffle",
			Parameter: "device",
		}
	}

	requestUrl, err := url.Parse(player.TransferUserPlayback)
	if err != nil {
		return err
	}

	transferPlayback := &object.TransferPlayback{
		Play:    play,
		Devices: []string{device.Id},
	}

	err = performNoResponsePut(accessToken, requestUrl.String(), transferPlayback)
	if err != nil {
		return err
	}

	return nil
}

func RecentlyPlayedTracksBefore(accessToken *accessToken.AccessToken, limit int, before time.Time) (*object.PlayItems, error) {
	requestUrl, err := url.Parse(player.GetCurrentUserRecentlyPlayedTracks)
	if err != nil {
		return nil, err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("limit", strconv.Itoa(limit))
	requestQuery.Add("before", strconv.FormatInt(before.Unix()*1000, 10))

	requestUrl.RawQuery = requestQuery.Encode()

	playHistory := new(object.PlayItems)
	err = performGet(accessToken, requestUrl.String(), playHistory)
	if err != nil {
		return nil, err
	}

	return playHistory, nil
}

func RecentlyPlayedTracksAfter(accessToken *accessToken.AccessToken, limit int, after time.Time) (*object.PlayItems, error) {
	requestUrl, err := url.Parse(player.GetCurrentUserRecentlyPlayedTracks)
	if err != nil {
		return nil, err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("limit", strconv.Itoa(limit))
	requestQuery.Add("after", strconv.FormatInt(after.Unix()*1000, 10))

	requestUrl.RawQuery = requestQuery.Encode()

	playHistory := new(object.PlayItems)
	err = performGet(accessToken, requestUrl.String(), playHistory)
	if err != nil {
		return nil, err
	}

	return playHistory, nil
}

func Play(accessToken *accessToken.AccessToken, device *object.Device) error {
	requestUrl, err := url.Parse(player.StartUserPlayback)
	if err != nil {
		return err
	}

	if device != nil {
		requestQuery := requestUrl.Query()
		requestQuery.Add("device_id", device.Id)

		requestUrl.RawQuery = requestQuery.Encode()
	}

	err = performEmptyNoResponsePut(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

const (
	RepeatModeTrack   object.RepeatMode = "track"
	RepeatModeContext object.RepeatMode = "context"
	RepeatModeOff     object.RepeatMode = "off"
)

func SetRepeatMode(accessToken *accessToken.AccessToken, state object.RepeatMode, device *object.Device) error {
	requestUrl, err := url.Parse(player.SetRepeatModeOnUserPlayback)
	if err != nil {
		return err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("state", string(state))
	if device != nil {
		requestQuery.Add("device_id", device.Id)
	}

	requestUrl.RawQuery = requestQuery.Encode()

	err = performEmptyNoResponsePut(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func GetCurrentlyPlaying(accessToken *accessToken.AccessToken, market string) (*object.CurrentlyPlaying, error) {
	requestUrl, err := url.Parse(player.GetInformationAboutTheUserCurrentPlayback)
	if err != nil {
		return nil, err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("market", market)

	requestUrl.RawQuery = requestQuery.Encode()

	currentlyPlaying := new(object.CurrentlyPlaying)
	err = performDebugGet(accessToken, requestUrl.String(), currentlyPlaying)
	if err != nil {
		return nil, err
	}

	return currentlyPlaying, nil
}

func GetCurrentlyPlayingTrack(accessToken *accessToken.AccessToken, market string) (*object.CurrentlyPlayingTrack, error) {
	requestUrl, err := url.Parse(player.GetUserCurrentlyPlayingTrack)
	if err != nil {
		return nil, err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("market", market)

	requestUrl.RawQuery = requestQuery.Encode()

	currentlyPlayingTrack := new(object.CurrentlyPlayingTrack)
	err = performDebugGet(accessToken, requestUrl.String(), currentlyPlayingTrack)
	if err != nil {
		return nil, err
	}

	return currentlyPlayingTrack, nil
}

func SetVolume(accessToken *accessToken.AccessToken, volumePercent int, device *object.Device) error {
	requestUrl, err := url.Parse(player.SetVolumeForUserPlayback)
	if err != nil {
		return err
	}

	requestQuery := requestUrl.Query()
	requestQuery.Add("volume_percent", strconv.Itoa(volumePercent))

	if device != nil {
		requestQuery.Add("device_id", device.Id)
	}

	requestUrl.RawQuery = requestQuery.Encode()

	err = performEmptyNoResponsePut(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func Pause(accessToken *accessToken.AccessToken, device *object.Device) error {
	requestUrl, err := url.Parse(player.PauseUserPlayback)
	if err != nil {
		return err
	}

	if device != nil {
		requestQuery := requestUrl.Query()
		requestQuery.Add("device_id", device.Id)

		requestUrl.RawQuery = requestQuery.Encode()
	}

	err = performEmptyNoResponsePut(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}

func PreviousTrack(accessToken *accessToken.AccessToken, device *object.Device) error {
	requestUrl, err := url.Parse(player.SkipUserPlaybackToPreviousTrack)
	if err != nil {
		return err
	}

	if device != nil {
		requestQuery := requestUrl.Query()
		requestQuery.Add("device_id", device.Id)

		requestUrl.RawQuery = requestQuery.Encode()
	}

	err = performEmptyNoResponsePost(accessToken, requestUrl.String())
	if err != nil {
		return err
	}

	return nil
}
