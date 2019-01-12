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

func Validate2xxResponse(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &UnexpectedResponseCodeError{
			Status:     response.Status,
			StatusCode: response.StatusCode,
		}
	}

	return nil
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

func Validate202Response(response *http.Response) error {
	if response.StatusCode != 202 {
		return &UnexpectedResponseCodeError{
			Status:     response.Status,
			StatusCode: response.StatusCode,
		}
	}

	return nil
}

func Validate204Response(response *http.Response) error {
	if response.StatusCode != 204 {
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

type RequiredParameterError struct {
	Method    string
	Parameter string
}

func (error *RequiredParameterError) Error() string {
	return fmt.Sprintf("RequiredParameter '%s' in method '%s'", error.Parameter, error.Method)
}
