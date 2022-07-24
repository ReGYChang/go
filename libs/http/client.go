package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"nexdata/pkg/config"
	"time"
)

const (
	Scheme        = "http"
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

func Request(method string, path string, payload io.Reader) (body []byte, err error) {
	// 如果 retry max 次數為0，直接送出
	if config.Destination.Retry.Max == 0 {
		return request(method, path, payload)
	}

	// retry，直到最大值
	for attempt := uint(0); attempt < config.Destination.Retry.Max; attempt++ {
		// 等候間隔時間
		if err := waitRetryBackoff(attempt); err != nil {
			return nil, err
		}

		body, err = request(method, path, payload)
		if err == nil {
			return body, err
		}
	}

	return body, err
}

func waitRetryBackoff(attempt uint) error {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = config.Destination.Retry.Interval
	}

	if waitTime > 0 {
		timer := time.NewTimer(waitTime)
		<-timer.C
	}

	return nil
}

// request request
func request(method string, path string, payload io.Reader) ([]byte, error) {
	u := &url.URL{
		Scheme: Scheme,
		Host:   config.Destination.HTTP.Host,
		Path:   path,
	}

	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
