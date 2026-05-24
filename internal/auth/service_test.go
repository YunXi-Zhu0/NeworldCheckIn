package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"checkin/internal/client"
)

func TestLoginRequest(t *testing.T) {
	var gotMethod string
	var gotContentType string
	var gotBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ret":1,"msg":"ok"}`))
	}))
	defer server.Close()

	httpClient, err := client.New()
	if err != nil {
		t.Fatal(err)
	}

	service := NewWithURL(httpClient, server.URL)
	result, err := service.Login("user@example.com", "secret")
	if err != nil {
		t.Fatal(err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q", gotMethod)
	}
	if gotContentType != "application/json" {
		t.Fatalf("content type = %q", gotContentType)
	}
	if gotBody != `{"email":"user@example.com","passwd":"secret"}` {
		t.Fatalf("body = %q", gotBody)
	}
	if result.Ret != 1 || result.Msg != "ok" {
		t.Fatalf("unexpected response: %+v", result)
	}
}
