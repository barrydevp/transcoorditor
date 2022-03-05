package util

import (
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func GetRequest() *resty.Client {
	return client
}

