package response

import (
	"fmt"
	"io"
)

const (
	CRLF        = "\r\n"
	HTTPVersion = "HTTP/1.1"
)

type StatusCode int

const (
	StatusResponseOK          StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 502
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	var reason string
	switch statusCode {
	case StatusResponseOK:
		reason = " OK"
	case StatusBadRequest:
		reason = " Bad Request"
	case StatusInternalServerError:
		reason = " Inernal Server Error"
	default:
		reason = ""

	}

	_, err := fmt.Fprintf(w, "%s %d%s%s", HTTPVersion, statusCode, reason, CRLF)
	return err
}
