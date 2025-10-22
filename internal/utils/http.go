package utils

import (
	"github.com/go-resty/resty/v2"
)

// Creates configured Resty client
func NewHTTPClient() *resty.Client {
	client := resty.New()
	//retries, timeouts, backoff
	return client
}
