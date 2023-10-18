package opensky

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	httpClient http.Client
	authHeader string
	baseUrl    string
}

func NewClient(opts ...ClientOption) (*Client, error) {
	cli := &Client{
		httpClient: http.Client{
			Timeout: 1 * time.Second,
		},
		baseUrl: "https://opensky-network.org/api",
	}

	for _, opt := range opts {
		if err := opt.apply(cli); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return cli, nil
}

func (cli *Client) GetStates() (StateResponse, error) {
	endpoint, _ := url.JoinPath(cli.baseUrl, "states/all")
	req, _ := http.NewRequest(http.MethodGet, endpoint, http.NoBody)

	req.URL.Query().Set("lamin", "50")
	req.URL.Query().Set("lamax", "50.5")
	req.URL.Query().Set("lomin", "3")
	req.URL.Query().Set("lomax", "3.5")
	var respObj StateResponse
	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return respObj, fmt.Errorf("failed to perform requeste: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var rawResp rawStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
		return respObj, fmt.Errorf("failed to decode response: %w", err)
	}

	respObj = rawResp.parse()

	return respObj, nil
}
