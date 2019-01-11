// GetAccessToken is based on Python library
// https://github.com/enriquegh/spotify-webplayer-token/blob/master/spotify_token.py
//

package accessToken

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/trejjam/spotify/spotifyError"
	"github.com/trejjam/spotify/url/login"
	"golang.org/x/net/publicsuffix"
)

type AccessToken struct {
	AccessToken string
	Expiration  int
}

var BonCookie = "MHwwfC0xODMxNzI2NTk2fC03NjkzMjUxNzAzMnwxfDF8MXwx"
var UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"

func initHttpClient() (*http.Client, error) {
	jar, err := cookiejar.New(
		&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		},
	)

	if err != nil {
		return nil, err
	}

	spotifyUrl := &url.URL{
		Scheme:     "https",
		Opaque:     "",
		Host:       login.Domain,
		Path:       login.Path,
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

	return &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}, nil
}

func getCsrfToken(client *http.Client) (string, error) {
	preLoginResponse, err := client.Get(login.Login)

	if err != nil {
		return "", err
	}

	if preLoginResponse.StatusCode != 200 {
		return "", &spotifyError.UnexpectedResponseCodeError{
			Status:     preLoginResponse.Status,
			StatusCode: preLoginResponse.StatusCode,
		}
	}

	cookies := client.Jar.Cookies(preLoginResponse.Request.URL)

	var token string
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if cookie.Name == "csrf_token" {
			token = cookie.Value
		}
	}

	if token == "" {
		return "", &spotifyError.EmptyCsrfError{
			Cookies: cookies,
		}
	}

	return token, nil
}

func processLogin(client *http.Client, username, password, csrfToken string) error {
	data := url.Values{}
	data.Add("remember", "")
	data.Add("username", username)
	data.Add("password", password)
	data.Add("csrf_token", csrfToken)

	loginResponse, err := client.PostForm(login.LoginApi, data)

	if err != nil {
		return err
	}

	if loginResponse.StatusCode != 200 {
		return &spotifyError.UnexpectedResponseCodeError{
			Status:     loginResponse.Status,
			StatusCode: loginResponse.StatusCode,
		}
	}

	return nil
}

func getAccessTokenUsingBrowse(client *http.Client) (*AccessToken, error) {
	browseRequest, err := http.NewRequest("GET", login.Browse, nil)
	if err != nil {
		return nil, err
	}

	browseRequest.Header.Set("user-agent", UserAgent)

	browseResponse, err := client.Do(browseRequest)
	if err != nil {
		return nil, err
	}

	if browseResponse.StatusCode != 200 {
		return nil, &spotifyError.UnexpectedResponseCodeError{
			Status:     browseResponse.Status,
			StatusCode: browseResponse.StatusCode,
		}
	}

	cookies := client.Jar.Cookies(browseResponse.Request.URL)

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

	if accessToken == "" {
		return nil, &spotifyError.EmptyAccessTokenError{
			Cookies: cookies,
		}
	}

	if expiration == 0 {
		return nil, &spotifyError.EmptyExpirationError{
			Cookies: cookies,
		}
	}

	return &AccessToken{
		AccessToken: accessToken,
		Expiration:  expiration,
	}, nil
}

func GetAccessToken(username string, password string) (*AccessToken, error) {
	client, err := initHttpClient()

	if err != nil {
		return nil, err
	}

	csrfToken, err := getCsrfToken(client)

	if err != nil {
		return nil, err
	}

	err = processLogin(client, username, password, csrfToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessTokenUsingBrowse(client)

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
