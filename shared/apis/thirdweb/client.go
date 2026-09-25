package thirdweb

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	baseUrl = "https://api.thirdweb.com/v1"

	endpointVerifyJWT = baseUrl + "/wallet/me"
)

type Client struct {
	client *http.Client
	apiKey string
	logger *zap.Logger
}

func NewClient(apiKey string, logger *zap.Logger) *Client {
	return &Client{
		client: &http.Client{},
		apiKey: apiKey,
		logger: logger,
	}
}

func (c *Client) VerifyJWT(jwt string) (*UserInfo, error) {
	return nil, nil
}
