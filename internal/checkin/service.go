package checkin

import (
	"bytes"
	"net/http"
	"net/url"

	"checkin/internal/client"
)

const DefaultCheckinURL = "https://neworld.space/user/checkin"

type AuthService interface {
	Login(email, passwd string) (client.APIResponse, error)
}

type HTTPClient interface {
	DoJSON(req *http.Request) (client.APIResponse, error)
}

type Result struct {
	Login   client.APIResponse
	Checkin client.APIResponse
}

type Service struct {
	authService AuthService
	httpClient  HTTPClient
	checkinURL  string
}

func New(authService AuthService, httpClient HTTPClient) *Service {
	return &Service{
		authService: authService,
		httpClient:  httpClient,
		checkinURL:  DefaultCheckinURL,
	}
}

func NewWithURL(authService AuthService, httpClient HTTPClient, checkinURL string) *Service {
	return &Service{
		authService: authService,
		httpClient:  httpClient,
		checkinURL:  checkinURL,
	}
}

func (s *Service) Run(email, passwd string) (Result, error) {
	loginResult, err := s.authService.Login(email, passwd)
	if err != nil {
		return Result{}, err
	}

	form := url.Values{}
	form.Set("checkin_type", "time")

	req, err := http.NewRequest(http.MethodPost, s.checkinURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	checkinResult, err := s.httpClient.DoJSON(req)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Login:   loginResult,
		Checkin: checkinResult,
	}, nil
}
