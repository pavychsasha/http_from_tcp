package response

import "fmt"

const (
	CRLF        = "\r\n"
	HTTPVersion = "HTTP/1.1"
)

type StatusCode int

const (
	StatusResponseOK          StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusBadGateway          StatusCode = 502
)

func GetStatusLine(statusCode StatusCode) string {
	switch statusCode {
	case StatusResponseOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""

	}
}

func GetWriterStatusLine(statusCode StatusCode) (res string) {
	reason := GetStatusLine(statusCode)

	if reason != "" {

		res = fmt.Sprintf("%s %d %s%s", HTTPVersion, statusCode, reason, CRLF)
	} else {
		res = fmt.Sprintf("%s %d%s", HTTPVersion, statusCode, CRLF)
	}
	return res

}
