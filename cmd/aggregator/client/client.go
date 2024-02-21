package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wvalencia19/tolling/types"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func (c *Client) AggregateInvoice(distance types.Distance) error {
	b, err := json.Marshal(distance)
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
