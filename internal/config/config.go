package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var ErrHelp = errors.New("help requested")

type Config struct {
	Email  string
	Passwd string
}

func Parse(args []string) (Config, error) {
	fs := flag.NewFlagSet("checkin", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := Config{}
	fs.StringVar(&cfg.Email, "email", os.Getenv("NEWORD_EMAIL"), "account email")
	fs.StringVar(&cfg.Passwd, "passwd", os.Getenv("NEWORD_PASSWD"), "account password")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return Config{}, ErrHelp
		}
		return Config{}, err
	}

	if fs.NArg() != 0 {
		return Config{}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	if cfg.Email == "" || cfg.Passwd == "" {
		return Config{}, errors.New("missing credentials: use --email/--passwd or NEWORD_EMAIL/NEWORD_PASSWD")
	}

	return cfg, nil
}
