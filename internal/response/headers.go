package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/pavychsasha/httpfromtcp/internal/headers"
)

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
