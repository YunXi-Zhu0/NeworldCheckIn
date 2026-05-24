package checkin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"checkin/internal/auth"
	"checkin/internal/client"
)

func TestRunReusesSession(t *testing.T) {
	var checkinCookie string
	var checkinBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/login":
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "abc123", Path: "/"})
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ret":1,"msg":"登录成功，欢迎回来"}`))
		case "/user/checkin":
			checkinCookie = r.Header.Get("Cookie")
			body, _ := io.ReadAll(r.Body)
			checkinBody = string(body)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ret":0,"msg":"今天已经签到过了"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	httpClient, err := client.New()
	if err != nil {
		t.Fatal(err)
	}

	authService := auth.NewWithURL(httpClient, server.URL+"/auth/login")
	service := NewWithURL(authService, httpClient, server.URL+"/user/checkin")

	result, err := service.Run("user@example.com", "secret")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(checkinCookie, "session=abc123") {
		t.Fatalf("cookie header = %q", checkinCookie)
	}
	if checkinBody != "checkin_type=time" {
		t.Fatalf("checkin body = %q", checkinBody)
	}
	if result.Login.Ret != 1 || result.Checkin.Ret != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
}
