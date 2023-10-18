package opensky

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// GetStates wraps GetStatesWithContext with context.Background().
func (cli *Client) GetStates(bbox *BoundingBox) (StateResponse, error) {
	return cli.GetStatesWithContext(context.Background(), bbox)
}

func (cli *Client) GetStatesWithContext(ctx context.Context, bbox *BoundingBox) (StateResponse, error) {
	endpoint, _ := url.JoinPath(cli.baseUrl, "states/all")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)

	if bbox != nil {
		q := req.URL.Query()

		q.Add("lamin", strconv.FormatFloat(bbox.LatitudeMin, 'f', -1, 32))
		q.Add("lamax", strconv.FormatFloat(bbox.LatitudeMax, 'f', -1, 32))
		q.Add("lomin", strconv.FormatFloat(bbox.LongitudeMin, 'f', -1, 32))
		q.Add("lomax", strconv.FormatFloat(bbox.LongitudeMax, 'f', -1, 32))

		req.URL.RawQuery = q.Encode()
	}

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
