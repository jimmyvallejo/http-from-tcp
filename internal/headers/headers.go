package headers

import (
	"bytes"
	"errors"
	"regexp"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

var headerNameRegex = regexp.MustCompile(`^[a-zA-Z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`)

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	contains := bytes.Contains(data, []byte("\r\n"))
	if !contains {
		return 0, false, nil
	}

	startsWith := bytes.HasPrefix(data, []byte("\r\n"))
	if startsWith {
		return 2, true, nil
	}

	trimmed := bytes.TrimSpace(data)

	splitted := bytes.SplitN(trimmed, []byte(":"), 2)

	toCheck := splitted[0]

	if !isValidHeaderName(toCheck) {
		return 0, false, errors.New("invalid header name")
	}

	fieldValue := bytes.TrimSpace(splitted[1])

	if bytes.ContainsAny(toCheck, " ") {
		return 0, false, errors.New("invalid format")
	}

	keyToLower := bytes.ToLower(toCheck)

	_, exists := h[string(keyToLower)]
	if exists {
		h[string(keyToLower)] += ", " + string(fieldValue)
	} else {
		h[string(keyToLower)] = string(fieldValue)
	}

	idx := bytes.Index(data, []byte("\r\n"))

	return idx + 2, false, nil

}

func isValidHeaderName(name []byte) bool {
	return headerNameRegex.Match(name)
}
