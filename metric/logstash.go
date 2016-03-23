// Copyright 2015 bs authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"encoding/json"
	"github.com/tsuru/bs/bslog"
	"net"
	"os"
)

func newLogStash() (statter, error) {
	const (
		defaultClient   = "tsuru"
		defaultPort     = "1984"
		defaultHost     = "localhost"
		defaultProtocol = "udp"
	)
	client := os.Getenv("METRICS_LOGSTASH_CLIENT")
	if client == "" {
		client = defaultClient
	}
	port := os.Getenv("METRICS_LOGSTASH_PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("METRICS_LOGSTASH_HOST")
	if host == "" {
		host = defaultHost
	}
	protocol := os.Getenv("METRICS_LOGSTASH_PROTOCOL")
	if protocol == "" {
		protocol = defaultProtocol
	}
	return &logStash{
		Client:   client,
		Host:     host,
		Port:     port,
		Protocol: protocol,
	}, nil
}

type logStash struct {
	Host     string
	Port     string
	Client   string
	Protocol string
}

func (s *logStash) Send(app, hostname, process, key string, value interface{}) error {
	message := map[string]interface{}{
		"client":  s.Client,
		"count":   1,
		"metric":  key,
		"value":   value,
		"app":     app,
		"host":    hostname,
		"process": process,
	}
	return s.send(message)
}

func (s *logStash) SendConn(app, hostname, process, host string) error {
	message := map[string]interface{}{
		"client":     s.Client,
		"count":      1,
		"metric":     "connection",
		"connection": host,
		"app":        app,
		"host":       hostname,
		"process":    process,
	}
	return s.send(message)
}

func (s *logStash) SendSys(hostname, key string, value interface{}) error {
	message := map[string]interface{}{
		"client": s.Client,
		"count":  1,
		"metric": key,
		"value":  value,
		"host":   hostname,
	}
	return s.send(message)
}

func (s *logStash) send(message map[string]interface{}) error {
	conn, err := net.Dial(s.Protocol, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		return err
	}
	defer conn.Close()
	data, err := json.Marshal(message)
	if err != nil {
		bslog.Errorf("unable to marshal metrics data json. Wrote %d bytes before error: %s", err)
		return err
	}
	bytesWritten, err := conn.Write(data)
	if err != nil {
		bslog.Errorf("unable to send metrics to logstash via UDP. Wrote %d bytes before error: %s", bytesWritten, err)
		return err
	}
	return nil
}
