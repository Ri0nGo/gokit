package oauth2Client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 3 * time.Second

	authorizePath = "/api/iam/oauth/authorize"
	tokenPath     = "/api/iam/oauth/token"
	refreshPath   = "/api/iam/oauth/refresh_token"
	authPath      = "/api/iam/oauth/auth"
	userinfoPath  = "/api/iam/oauth/userinfo"
)

// Config contains OAuth2 client settings for the IAM service.
type Config struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	Timeout      time.Duration
}

// Client calls the IAM OAuth2 authorization-code endpoints.
type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	redirectURI  string
	scopes       []string
	httpClient   *http.Client
}

// TokenResponse is returned by the IAM token endpoint.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope,omitempty"`
}

// UserInfo is returned by the IAM userinfo endpoint.
type UserInfo struct {
	OpenID      string   `json:"openid"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Status      int      `json:"status"`
	Roles       []string `json:"roles"`
}

// New creates an IAM OAuth2 client.
func New(config Config) (*Client, error) {
	baseURLValue := strings.TrimSpace(config.BaseURL)
	clientID := strings.TrimSpace(config.ClientID)
	clientSecret := strings.TrimSpace(config.ClientSecret)
	redirectURI := strings.TrimSpace(config.RedirectURI)
	if baseURLValue == "" {
		return nil, errors.New("base url is required")
	}
	if clientID == "" {
		return nil, errors.New("client id is required")
	}
	if clientSecret == "" {
		return nil, errors.New("client secret is required")
	}
	if redirectURI == "" {
		return nil, errors.New("redirect uri is required")
	}
	baseURL, err := url.Parse(strings.TrimRight(baseURLValue, "/"))
	if err != nil || baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("invalid base url: %s", config.BaseURL)
	}
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &Client{
		baseURL:      baseURL.String(),
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		scopes:       append([]string(nil), config.Scopes...),
		httpClient:   &http.Client{Timeout: timeout},
	}, nil
}

// AuthCodeURL returns the IAM authorization URL that the user should visit.
func (c *Client) AuthCodeURL(state string, scopes ...string) (string, error) {
	endpoint, err := url.Parse(c.baseURL + authorizePath)
	if err != nil {
		return "", err
	}
	if len(scopes) == 0 {
		scopes = c.scopes
	}
	query := endpoint.Query()
	query.Set("response_type", "code")
	query.Set("client_id", c.clientID)
	query.Set("redirect_uri", c.redirectURI)
	if state != "" {
		query.Set("state", state)
	}
	if len(scopes) > 0 {
		query.Set("scope", strings.Join(scopes, " "))
	}
	endpoint.RawQuery = query.Encode()
	return endpoint.String(), nil
}

// GetUserToken exchanges an authorization code for a user access token.
func (c *Client) GetUserToken(ctx context.Context, code string) (*TokenResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("code is required")
	}
	endpoint, err := c.endpoint(tokenPath, map[string]string{
		"client_id":  c.clientID,
		"secret":     c.clientSecret,
		"code":       strings.TrimSpace(code),
		"grant_type": "authorization_code",
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	var token TokenResponse
	if err := c.doAPI(req, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// RefreshUserToken refreshes an OAuth2 access token with a refresh token.
func (c *Client) RefreshUserToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, errors.New("refresh token is required")
	}
	endpoint, err := c.endpoint(refreshPath, map[string]string{
		"client_id":     c.clientID,
		"grant_type":    "refresh_token",
		"refresh_token": strings.TrimSpace(refreshToken),
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	var token TokenResponse
	if err := c.doAPI(req, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// CheckUserToken verifies that an access token is valid for the openid.
func (c *Client) CheckUserToken(ctx context.Context, accessToken string, openID string) error {
	if strings.TrimSpace(accessToken) == "" {
		return errors.New("access token is required")
	}
	if strings.TrimSpace(openID) == "" {
		return errors.New("openid is required")
	}
	endpoint, err := c.endpoint(authPath, map[string]string{
		"access_token": strings.TrimSpace(accessToken),
		"openid":       strings.TrimSpace(openID),
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	var response struct {
		Success bool `json:"success"`
	}
	if err := c.doAPI(req, &response); err != nil {
		return err
	}
	if !response.Success {
		return errors.New("iam token check failed")
	}
	return nil
}

// GetUserInfo returns the current user bound to the access token.
func (c *Client) GetUserInfo(ctx context.Context, accessToken string, openID string) (*UserInfo, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("access token is required")
	}
	if strings.TrimSpace(openID) == "" {
		return nil, errors.New("openid is required")
	}
	endpoint, err := c.endpoint(userinfoPath, map[string]string{
		"access_token": strings.TrimSpace(accessToken),
		"openid":       strings.TrimSpace(openID),
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	var user UserInfo
	if err := c.doAPI(req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) endpoint(path string, values map[string]string) (string, error) {
	endpoint, err := url.Parse(c.baseURL + path)
	if err != nil {
		return "", err
	}
	query := endpoint.Query()
	for key, value := range values {
		query.Set(key, value)
	}
	endpoint.RawQuery = query.Encode()
	return endpoint.String(), nil
}

func (c *Client) doJSON(req *http.Request, target any) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return parseError(res.StatusCode, body)
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode iam response: %w", err)
	}
	return nil
}

func (c *Client) doAPI(req *http.Request, target any) error {
	var response apiResponse
	if err := c.doJSON(req, &response); err != nil {
		return err
	}
	if response.Code != 0 {
		return fmt.Errorf("iam api error: %s", response.message())
	}
	if target == nil {
		return nil
	}
	if len(response.Data) == 0 {
		return errors.New("iam api response missing data")
	}
	if err := json.Unmarshal(response.Data, target); err != nil {
		return fmt.Errorf("decode iam data: %w", err)
	}
	return nil
}

func parseError(statusCode int, body []byte) error {
	var apiErr apiResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.message() != "" {
		return fmt.Errorf("iam api error: %s", apiErr.message())
	}
	return fmt.Errorf("iam request failed with status %d", statusCode)
}

type apiResponse struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (r apiResponse) message() string {
	if r.Msg != "" {
		return r.Msg
	}
	return r.Message
}
