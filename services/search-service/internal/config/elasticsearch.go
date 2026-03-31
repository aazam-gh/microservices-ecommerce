package config

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearch(cfg *Config) (*elasticsearch.Client, error) {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DialContext:         (&net.Dialer{Timeout: cfg.Elastic.Timeout}).DialContext,
		TLSHandshakeTimeout: cfg.Elastic.Timeout,
	}

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Elastic.URL},
		Transport: transport,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
