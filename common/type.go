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
	"Date",
	"Content-Type",
	"Content-Length",
	"Content-Encoding",
	"ccess-Control-Allow-Origin",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Methods",
	"Access-Control-Expose-Headers",
}

type RequestPackage struct {
	URL    *url.URL
	Body   []byte
	Method string
	Header http.Header
}

type ResponsePackage struct {
	StatusCode int
	Body       []byte
	HeadMap    map[string]string
}
