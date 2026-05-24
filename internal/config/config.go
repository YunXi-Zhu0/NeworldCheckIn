package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrHelp = errors.New("help requested")

const defaultConfigPath = "config.yaml"

type Config struct {
	Email  string
	Passwd string
}

func Parse(args []string) (Config, error) {
	fs := flag.NewFlagSet("checkin", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg, err := loadConfig(defaultConfigPath)
	if err != nil {
		return Config{}, err
	}

	if envEmail := os.Getenv("NEWORD_EMAIL"); envEmail != "" {
		cfg.Email = envEmail
	}
	if envPasswd := os.Getenv("NEWORD_PASSWD"); envPasswd != "" {
		cfg.Passwd = envPasswd
	}

	fs.StringVar(&cfg.Email, "email", cfg.Email, "account email")
	fs.StringVar(&cfg.Passwd, "passwd", cfg.Passwd, "account password")

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
		return Config{}, errors.New("missing credentials: use --email/--passwd, NEWORD_EMAIL/NEWORD_PASSWD, or config.yaml")
	}

	return cfg, nil
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config.yaml: %w", err)
	}

	return cfg, nil
}
