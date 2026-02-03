package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

const CRLF = "\r\n"

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
	h[key] = value
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

	key = bytes.ReplaceAll(key, []byte(" "), []byte(""))
	value = bytes.ReplaceAll(value, []byte(" "), []byte(""))
	return key, value, nil
}
