package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/pavychsasha/httpfromtcp/internal/headers"
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

func GetDefaultHeaders(contentLen int) headers.Headers {

	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s:%s%s", key, value, CRLF)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "%s", CRLF)
	if err != nil {
		return err
	}
	return nil
}
