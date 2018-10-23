package common

import (
	"net/http"
	"net/url"
)

const (
	PROTOCOL = "https"
	PKI_PATH = "/etc/kubernetes/pki/"
)

var HEAD_LIST = []string{
	"Content-Type",
	"Date",
	"Content-Length",
}

type RequestPackage struct {
	URL    *url.URL
	Body   []byte
	Method string
	Header http.Header
}

type ResponsePackage struct {
	Body       []byte
	HeadMap    map[string]string
	StatusCode int
}
