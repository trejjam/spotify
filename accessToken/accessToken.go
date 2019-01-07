// GetAccessToken is based on Python library
// https://github.com/enriquegh/spotify-webplayer-token/blob/master/spotify_token.py
//

package accessToken

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"golang.org/x/net/publicsuffix"
)

type AccessToken struct {
	AccessToken string
	Expiration  int
}

var BonCookie = "Eeg8phaiKah4JeekirooGhoh3faehoo0yieghoeghahlaeX7a"
var UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"

var spotifyLoginUrl = "https://accounts.spotify.com/login"
var spotifyLoginApiUrl = "https://accounts.spotify.com/api/login"
var spotifyBrowseUrl = "https://open.spotify.com/browse"

func GetAccessToken(username string, password string) AccessToken {
	jar, err := cookiejar.New(
		&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		},
	)

	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	spotifyUrl := &url.URL{
		Scheme:     "https",
		Opaque:     "",
		Host:       "accounts.spotify.com",
		Path:       "/login",
		ForceQuery: false,
		RawPath:    "",
		Fragment:   "",
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "__bon",
		Value:  BonCookie,
		Path:   "/",
		Domain: spotifyUrl.Host,
	}
	cookies = append(cookies, cookie)

	jar.SetCookies(spotifyUrl, cookies)

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}
	preLoginResponse, err := client.Get(spotifyLoginUrl)

	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	if preLoginResponse.StatusCode != 200 {
		fmt.Errorf(preLoginResponse.Status)
		panic(preLoginResponse.Status)
	}

	cookies = client.Jar.Cookies(preLoginResponse.Request.URL)

	var token string
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if cookie.Name == "csrf_token" {
			token = cookie.Value
		}
	}

	data := url.Values{}
	data.Add("remember", "")
	data.Add("username", username)
	data.Add("password", password)
	data.Add("csrf_token", token)

	loginResponse, err := client.PostForm(spotifyLoginApiUrl, data)

	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	if loginResponse.StatusCode != 200 {
		fmt.Errorf(loginResponse.Status)
		panic(loginResponse.Status)
	}

	browseRequest, err := http.NewRequest("GET", spotifyBrowseUrl, nil)
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	browseRequest.Header.Set("user-agent", UserAgent)

	browseResponse, err := client.Do(browseRequest)
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}
	if browseResponse.StatusCode != 200 {
		fmt.Errorf(browseResponse.Status)
		panic(browseResponse.Status)
	}

	cookies = client.Jar.Cookies(browseResponse.Request.URL)

	var accessToken string
	var expiration int
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]

		switch cookie.Name {
		case "wp_access_token":
			accessToken = cookie.Value
		case "wp_expiration":
			expiration, err = strconv.Atoi(cookie.Value)
			expiration /= 1000
		}
	}

	return AccessToken{
		AccessToken: accessToken,
		Expiration:  expiration,
	}
}
