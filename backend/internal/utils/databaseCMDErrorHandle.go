package utils

import (
	"bytes"
)

func ExtractErrorMessage(stderr bytes.Buffer) string {
	return stderr.String()
}
