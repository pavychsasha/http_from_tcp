package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	crlfIdx := bytes.IndexAny(data, CRLF)
	switch crlfIdx {
	case -1:
		return 0, false, nil
	case 0:
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	key, value, err := parseHeaderBytes(data[:crlfIdx])
	if err != nil {
		return 0, false, fmt.Errorf("Could not parse headers string: %v", err)
	}

	h.Set(string(key), string(value))
	return crlfIdx + 2, false, nil

}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func parseHeaderBytes(val []byte) (key []byte, value []byte, err error) {

	values := bytes.SplitN(val, []byte(":"), 2)
	if len(values) < 2 {
		return []byte{}, []byte{}, fmt.Errorf("Invalid header format: headers should be in the following format:\nfield-line   = field-name ':' OWS field-value OWS\nGiven: %s", val)
	}
	key = values[0]
	value = values[1]

	if key[len(key)-1] == byte(' ') {
		return []byte{}, []byte{}, fmt.Errorf("Invalid header: there should not be any spaces between field-name and colon in %s", val)
	}
	key = bytes.ReplaceAll(bytes.ToLower(key), []byte(" "), []byte(""))
	value = bytes.ReplaceAll(value, []byte(" "), []byte(""))

	if !validTokens(key) {
		return []byte{}, []byte{}, fmt.Errorf("Invalid header token found: %s", key)
	}

	return key, value, nil
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}
