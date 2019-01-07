package spotify

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
		Value:  "MHwwfC01ODc4MjExMzJ8LTI0Njg4NDg3NTQ0fDF8MXwxfDE",
		Path:   "/",
		Domain: spotifyUrl.Host,
	}
	cookies = append(cookies, cookie)

	jar.SetCookies(spotifyUrl, cookies)

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}
	preLoginResponse, err := client.Get("https://accounts.spotify.com/login")

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

	loginResponse, err := client.PostForm("https://accounts.spotify.com/api/login", data)

	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	if loginResponse.StatusCode != 200 {
		fmt.Errorf(loginResponse.Status)
		panic(loginResponse.Status)
	}

	browseRequest, err := http.NewRequest("GET", "https://open.spotify.com/browse", nil)
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}

	browseRequest.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")

	browseResponse, err := client.Do(browseRequest)
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}
	if browseResponse.StatusCode != 200 {
		fmt.Errorf(browseResponse.Status)
		panic(browseResponse.Status)
	}

	fmt.Println(browseResponse.Status)

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
