package app

import (
	"errors"
	"fmt"
	"io"

	"checkin/internal/auth"
	"checkin/internal/checkin"
	"checkin/internal/client"
	"checkin/internal/config"
	"checkin/internal/output"
)

type runner interface {
	Run(email, passwd string) (checkin.Result, error)
}

var newRunner = func() (runner, error) {
	httpClient, err := client.New()
	if err != nil {
		return nil, err
	}

	authService := auth.New(httpClient)
	return checkin.New(authService, httpClient), nil
}

func Run(args []string, stdout, stderr io.Writer) int {
	cfg, err := config.Parse(args)
	if err != nil {
		if errors.Is(err, config.ErrHelp) {
			output.WriteUsage(stdout)
			return 0
		}

		fmt.Fprintln(stderr, err)
		output.WriteUsage(stderr)
		return 2
	}

	checkinRunner, err := newRunner()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	result, err := checkinRunner.Run(cfg.Email, cfg.Passwd)
	if err != nil {
		output.WriteError(stderr, err)
		return 1
	}

	if err := output.WriteResult(stdout, result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	return 0
}
