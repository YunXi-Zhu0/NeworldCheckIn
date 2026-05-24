package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"checkin/internal/checkin"
	"checkin/internal/client"
)

func WriteResult(w io.Writer, result checkin.Result) error {
	if _, err := fmt.Fprintln(w, "login:"); err != nil {
		return err
	}
	if err := writeJSON(w, result.Login); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "checkin:"); err != nil {
		return err
	}
	return writeJSON(w, result.Checkin)
}

func WriteError(w io.Writer, err error) {
	var httpErr client.HTTPError
	if errors.As(err, &httpErr) {
		fmt.Fprintln(w, httpErr.Error())
		if httpErr.Body != "" {
			fmt.Fprintln(w, httpErr.Body)
		}
		return
	}

	fmt.Fprintf(w, "request failed: %v\n", err)
}

func WriteUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  checkin [--email EMAIL] [--passwd PASSWD]")
}

func writeJSON(w io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(data))
	return err
}
