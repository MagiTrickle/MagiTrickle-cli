package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	api "github.com/Ponywka/MagiTrickle/backend/pkg/api"
)

func doUnixRequest(method, urlPath string, body []byte) (*http.Response, error) {
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", api.SocketPath)
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, "http://unix"+urlPath, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request fail: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func doUnixJSON(method, urlPath string, data interface{}) (*http.Response, error) {
	var body []byte
	if data != nil {
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("json marshal fail: %w", err)
		}
	}
	return doUnixRequest(method, urlPath, body)
}
