package common

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpSend struct {
	Host      string
	CertPool  *x509.CertPool
	ClientCrt tls.Certificate
}

type HttpSender interface {
	SendHttpRequest(requestPackage RequestPackage) (ResponsePackage, error)
}

func NewHttpSend(dstAddr string) (HttpSender, error) {
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(PKI_PATH + "ca.crt")
	if err != nil {
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)
	cliCrt, err := tls.LoadX509KeyPair(PKI_PATH+"client.crt", PKI_PATH+"client.key")
	if err != nil {
		return nil, err
	}
	return &httpSend{
		Host:      fmt.Sprintf("%s://%s", PROTOCOL, dstAddr),
		CertPool:  pool,
		ClientCrt: cliCrt,
	}, nil
}

func (self httpSend) SendHttpRequest(requestPackage RequestPackage) (ResponsePackage, error) {
	var responsePackage ResponsePackage
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      self.CertPool,
				Certificates: []tls.Certificate{self.ClientCrt},
			},
		},
	}

	body := ioutil.NopCloser(strings.NewReader(string(requestPackage.Body)))
	req, err := http.NewRequest(requestPackage.Method, fmt.Sprintf("%s%s", self.Host, requestPackage.URL), body)
	if err != nil {
		return responsePackage, err
	}
	defer req.Body.Close()
	req.Header = requestPackage.Header

	respon, err := client.Do(req)
	if err != nil {
		return responsePackage, err
	}
	defer respon.Body.Close()

	headMap := make(map[string]string, len(HEAD_LIST))
	for _, key := range HEAD_LIST {
		if value := respon.Header.Get(key); value != "" {
			headMap[key] = value
		}
	}

	responsePackage.StatusCode = respon.StatusCode
	responsePackage.HeadMap = headMap
	responsePackage.Body, _ = ioutil.ReadAll(respon.Body)
	return responsePackage, nil
}
