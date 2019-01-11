package spotifyError

import (
	"fmt"
	"net/http"
)

type UnexpectedResponseCodeError struct {
	Status     string
	StatusCode int
}

func (error *UnexpectedResponseCodeError) Error() string {
	return fmt.Sprintf("Unexpected response code %s", error.Status)
}

func Validate200Response(response *http.Response) error {
	if response.StatusCode != 200 {
		return &UnexpectedResponseCodeError{
			Status:     response.Status,
			StatusCode: response.StatusCode,
		}
	}

	return nil
}

type EmptyCsrfError struct {
	Cookies []*http.Cookie
}

func (error *EmptyCsrfError) Error() string {
	return "EmptyCsrfError"
}

type EmptyAccessTokenError struct {
	Cookies []*http.Cookie
}

func (error *EmptyAccessTokenError) Error() string {
	return "EmptyAccessTokenError"
}

type EmptyExpirationError struct {
	Cookies []*http.Cookie
}

func (error *EmptyExpirationError) Error() string {
	return "EmptyExpirationError"
}
