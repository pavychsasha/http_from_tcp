package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pavychsasha/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine  RequestLine
	RequestState RequestState
	Headers      headers.Headers
	// Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestState int

const (
	INITIALIZED RequestState = iota
	ParsingHeaders
	DONE
)

const CRLF = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	request := Request{
		RequestState: INITIALIZED,
		Headers:      headers.NewHeaders(),
	}

	for request.RequestState != DONE {
		if len(buffer) <= readToIndex {
			dst := make([]byte, len(buffer)*2)
			copy(dst, buffer)
			buffer = dst
		}

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if request.RequestState == INITIALIZED {
					return nil, fmt.Errorf("incomplete request")
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	crlfIndex := bytes.Index(data, []byte(CRLF))

	if crlfIndex == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:crlfIndex])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, crlfIndex + len(CRLF), nil
}

func requestLineFromString(str string) (*RequestLine, error) {

	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if 'A' > c || 'Z' < c {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}
	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		HttpVersion:   version,
		RequestTarget: parts[1],
	}, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.RequestState {
	case INITIALIZED:
		requestLine, numBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numBytes == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.RequestState = ParsingHeaders
		return numBytes, nil

	case ParsingHeaders:
		numBytes, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.RequestState = DONE
		}

		return numBytes, err
	case DONE:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func (r *Request) parse(data []byte) (int, error) {

	totalProcessedBytes := 0

	// Header line can contain multiple lines
	// Making sure that we are setting a DONE state after
	// processing all of the headers
	for r.RequestState != DONE {
		n, err := r.parseSingle(data[totalProcessedBytes:])
		if err != nil {
			return totalProcessedBytes, err
		}
		totalProcessedBytes += n
		if n == 0 {

			break
		}
	}
	return totalProcessedBytes, nil
}
