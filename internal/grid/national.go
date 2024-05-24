package grid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mailru/easyjson"
)

var (
	ErrEmptyData  = errors.New("empty data in response")
	timeFormat    = "2006-01-02T15:04:05Z"
	intensityPath = "/intensity"
)

type HTTPWrapper interface {
	MakeGetRequest(ctx context.Context, url string) ([]byte, int, error)
}

// Client is responsible for interaction with external API.
type Client struct {
	intensityURL string
	HTTPWrapper  HTTPWrapper
}

func NewClient(Domain string, HTTPWrapper HTTPWrapper) *Client {
	return &Client{intensityURL: Domain + intensityPath, HTTPWrapper: HTTPWrapper}
}

func (c *Client) GetIntensity(ctx context.Context) (val Intensity, err error) {
	// added time.Second, since api returns "14-14:30" for 14:30 time.
	// Can be replaced with */31,1 * * * *, but in that case update will be with one minute delay.
	data, status, err := c.HTTPWrapper.MakeGetRequest(ctx, fmt.Sprintf("%s/%s", c.intensityURL, time.Now().Add(time.Second).UTC().Format(timeFormat)))
	if err != nil {
		return "", fmt.Errorf("client call: make request: %v", err)
	}

	if status != http.StatusOK {
		return "", fmt.Errorf("status in not 200")
	}

	var r Response
	err = easyjson.Unmarshal(data, &r)
	if err != nil {
		return "", err
	}

	if len(r.Data) != 0 {
		return r.Data[0].Intensity.Index, nil
	}

	return "", ErrEmptyData
}
