package auth

import (
	"bytes"
	"encoding/json"
	"net/http"

	"checkin/internal/client"
)

const DefaultLoginURL = "https://neworld.space/auth/login"

type HTTPClient interface {
	DoJSON(req *http.Request) (client.APIResponse, error)
}

type Service struct {
	httpClient HTTPClient
	loginURL   string
}

func New(httpClient HTTPClient) *Service {
	return &Service{
		httpClient: httpClient,
		loginURL:   DefaultLoginURL,
	}
}

func NewWithURL(httpClient HTTPClient, loginURL string) *Service {
	return &Service{
		httpClient: httpClient,
		loginURL:   loginURL,
	}
}

func (s *Service) Login(email, passwd string) (client.APIResponse, error) {
	payload, err := json.Marshal(map[string]string{
		"email":  email,
		"passwd": passwd,
	})
	if err != nil {
		return client.APIResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, s.loginURL, bytes.NewReader(payload))
	if err != nil {
		return client.APIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	return s.httpClient.DoJSON(req)
}
