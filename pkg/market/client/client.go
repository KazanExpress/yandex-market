package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// DefaultAPIEndpoint is a default yandex market api endpoint.
const DefaultAPIEndpoint = "https://api.partner.market.yandex.ru/"

// YandexMarketClient wraps API calls to yandex market.
type YandexMarketClient struct {
	options *Options
}

// Options client constructor params.
type Options struct {
	OAuthToken    string
	OAuthClientID string
	APIEndpoint   string
	UserAgent     string
	Client        *http.Client
	Logger        *zap.Logger
}

// Option modifies Options.
type Option func(*Options)

// WithOAuth configures oauth clientID and token.
func WithOAuth(token, clientID string) Option {
	return func(o *Options) {
		o.OAuthClientID = clientID
		o.OAuthToken = token
	}
}

// WithLogger configures logger.
func WithLogger(logger *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithHTTPClient configures http client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *Options) {
		o.Client = client
	}
}

// WithUserAgent sets useragent.
func WithUserAgent(useragent string) Option {
	return func(o *Options) {
		o.UserAgent = useragent
	}
}

// NewYandexMarketClient is YandexMarketClient constructor.
func NewYandexMarketClient(opts ...Option) *YandexMarketClient {
	opt := &Options{
		Client:      http.DefaultClient,
		APIEndpoint: DefaultAPIEndpoint,
		Logger:      zap.NewNop(),
		UserAgent:   "KE/yandex-market client github.com/KazanExpress/yandex-market",
	}

	for _, o := range opts {
		o(opt)
	}

	return &YandexMarketClient{
		options: opt,
	}
}

func (c *YandexMarketClient) newRequest(
	ctx context.Context,
	method, path, query string,
	body io.Reader,
) (*http.Request, error) {
	fullURL := c.options.APIEndpoint
	// safe concat of API endpoint and path
	if !strings.HasSuffix(fullURL, "/") {
		fullURL += "/"
	}

	path = strings.TrimPrefix(path, "/")
	fullURL += path

	if !strings.HasSuffix(fullURL, ".json") {
		fullURL += ".json"
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("authorization",
		fmt.Sprintf("OAuth oauth_token=%s, oauth_client_id=%s",
			c.options.OAuthToken, c.options.OAuthClientID))
	req.Header.Add("user-agent", c.options.UserAgent)
	req.Header.Add("accept", "*/*")

	req.URL.RawQuery = query

	return req, nil
}

func (c *YandexMarketClient) executeRequest(req *http.Request, jsonResponse interface{}) error {
	resp, err := c.options.Client.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(jsonResponse); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.options.Logger.Error("failed to close response body",
				zap.Error(err),
			)
		}
	}()

	return nil
}
