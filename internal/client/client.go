package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

type APIResponse struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http error: %d", e.StatusCode)
}

type Client struct {
	httpClient *http.Client
}

func New() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: &http.Client{Jar: jar},
	}, nil
}

func (c *Client) DoJSON(req *http.Request) (APIResponse, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return APIResponse{}, HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return APIResponse{}, err
	}

	return result, nil
}
