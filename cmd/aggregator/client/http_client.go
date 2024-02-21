package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wvalencia19/tolling/types"
)

type HTTPClient struct {
	endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("the server responded with a non 200 status code %s", res.Status)
	}

	return nil
}
