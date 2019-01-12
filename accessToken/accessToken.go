// GetAccessToken is based on Python library
// https://github.com/enriquegh/spotify-webplayer-token/blob/master/spotify_token.py
//

package accessToken

import (
	"fmt"
	"io"
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

func setHeaders(request *http.Request, accessToken *AccessToken) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.AccessToken))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
}

func (accessToken *AccessToken) CreateGetRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setHeaders(request, accessToken)

	return request, nil
}

func (accessToken *AccessToken) CreatePostRequest(url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	setHeaders(request, accessToken)

	return request, nil
}

func (accessToken *AccessToken) CreatePutRequest(url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	setHeaders(request, accessToken)

	return request, nil
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

	spotifyUrl, err := url.Parse(login.Login)

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

	err = spotifyError.Validate200Response(preLoginResponse)
	if err != nil {
		return "", err
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

	err = spotifyError.Validate200Response(loginResponse)
	if err != nil {
		return err
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

	err = spotifyError.Validate200Response(browseResponse)
	if err != nil {
		return nil, err
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
