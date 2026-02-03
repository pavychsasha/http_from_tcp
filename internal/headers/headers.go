package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	crlfIdx := bytes.IndexAny(data, CRLF)
	switch crlfIdx {
	case -1:
		return 0, false, nil
	case 0:
		return 1, true, nil
	}

	headersFull := string(data[:crlfIdx])
	key, value, err := parseHeaderString(headersFull)
	if err != nil {
		return 0, false, fmt.Errorf("Could not parse headers string: %v", err)
	}
	h[key] = value

	return crlfIdx + 2, false, nil

}

func parseHeaderString(val string) (key string, value string, err error) {

	values := strings.SplitN(val, ":", 2)
	if len(values) < 2 {
		return "", "", fmt.Errorf("Invalid header format: headers should be in the following format:\nfield-line   = field-name ':' OWS field-value OWS\nGiven: %s", val)
	}
	key = values[0]
	value = values[1]

	if string(key[len(key)-1]) == " " {
		return "", "", fmt.Errorf("Invalid header: there should not be any spaces between field-name and colon in %s", val)
	}
	fmt.Println(key, " ", value)

	key = strings.ReplaceAll(key, " ", "")
	value = strings.ReplaceAll(value, " ", "")
	return key, value, nil
}
