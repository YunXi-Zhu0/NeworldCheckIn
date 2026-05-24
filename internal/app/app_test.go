package app

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"checkin/internal/checkin"
	"checkin/internal/client"
)

func TestRunOutputsSections(t *testing.T) {
	oldRunner := newRunner
	defer func() { newRunner = oldRunner }()

	newRunner = func() (runner, error) {
		return stubRunner{
			runFn: func(email, passwd string) (checkin.Result, error) {
				if email != "user@example.com" || passwd != "secret" {
					t.Fatalf("unexpected credentials: %s %s", email, passwd)
				}
				return checkin.Result{
					Login:   client.APIResponse{Ret: 1, Msg: "登录成功，欢迎回来"},
					Checkin: client.APIResponse{Ret: 0, Msg: "今天已经签到过了"},
				}, nil
			},
		}, nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--email", "user@example.com", "--passwd", "secret"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d stderr = %q", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "login:\n{\n  \"ret\": 1,\n  \"msg\": \"登录成功，欢迎回来\"\n}\n") {
		t.Fatalf("missing login section: %q", out)
	}
	if !strings.Contains(out, "checkin:\n{\n  \"ret\": 0,\n  \"msg\": \"今天已经签到过了\"\n}\n") {
		t.Fatalf("missing checkin section: %q", out)
	}
}

func TestRunMissingCredentials(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "")
	t.Setenv("NEWORD_PASSWD", "")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stderr.String(), "missing credentials") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stdout.String(), "checkin [--email EMAIL] [--passwd PASSWD]") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunHTTPError(t *testing.T) {
	oldRunner := newRunner
	defer func() { newRunner = oldRunner }()

	newRunner = func() (runner, error) {
		return stubRunner{
			runFn: func(email, passwd string) (checkin.Result, error) {
				return checkin.Result{}, client.HTTPError{
					StatusCode: 500,
					Body:       `{"ret":0,"msg":"server error"}`,
				}
			},
		}, nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--email", "user@example.com", "--passwd", "secret"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stderr.String(), "http error: 500") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunRequestError(t *testing.T) {
	oldRunner := newRunner
	defer func() { newRunner = oldRunner }()

	newRunner = func() (runner, error) {
		return stubRunner{
			runFn: func(email, passwd string) (checkin.Result, error) {
				return checkin.Result{}, errors.New("boom")
			},
		}, nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--email", "user@example.com", "--passwd", "secret"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code = %d", code)
	}
	if !strings.Contains(stderr.String(), "request failed: boom") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

type stubRunner struct {
	runFn func(email, passwd string) (checkin.Result, error)
}

func (s stubRunner) Run(email, passwd string) (checkin.Result, error) {
	return s.runFn(email, passwd)
}
