package pkg

import (
	"io/ioutil"
	"net/http"
	"time"
)

func GetBody(url string) ([]byte, error) {
	res := &http.Response{}
	var err error

	for {
		res, err = http.Get(url)
		if err != nil {
			return nil, err
		}

		switch {
		case isSuccess(res.StatusCode):
			break

		case isTooManyReq(res.StatusCode):
			time.Sleep(20 * time.Second)
			continue

		default:
			return nil, err
		}
		break
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

func isTooManyReq(code int) bool {
	return code == 429
}
