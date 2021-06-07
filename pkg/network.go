package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	err = handleStatus(res.StatusCode)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return body, err
	}

	return body, nil
}

func handleStatus(statusCode int) error {
	switch {
	case isSuccess(statusCode):
		return nil

	case isClientError(statusCode):
		return fmt.Errorf("Response status: %d ", statusCode)

	case isServerError(statusCode):
		return fmt.Errorf("Response status: %d ", statusCode)

	default:
		return nil
	}
}

func isSuccess(code int) bool {
	switch code {
	case 200:
		return true
	case 201:
		return true
	default:
		return false
	}
}

func isClientError(code int) bool {
	return code >= 400 && code < 500
}

func isServerError(code int) bool {
	return code >= 500
}
