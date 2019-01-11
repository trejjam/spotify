package player

import (
	"encoding/json"
	"github.com/trejjam/spotify/accessToken"
	"github.com/trejjam/spotify/object"
	"github.com/trejjam/spotify/spotifyError"
	"github.com/trejjam/spotify/url/player"
	"net/http"
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

func GetDevices(accessToken *accessToken.AccessToken) (*object.Devices, error) {
	request, err := accessToken.CreateGetRequest(player.GetUserAvailableDevices)
	if err != nil {
		return nil, err
	}

	client := getClient()

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	err = spotifyError.Validate200Response(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	devices := new(object.Devices)
	err = json.NewDecoder(response.Body).Decode(devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
