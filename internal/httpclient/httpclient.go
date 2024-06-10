package httpclient

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

func NewHttpClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
}
