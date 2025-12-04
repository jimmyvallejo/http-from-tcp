package request

import (
	"errors"
	"strconv"
	"strings"
)

func (r *Request) parseBody(data []byte) (int, error) {

	length, ok := r.Headers.Get("Content-Length")
	if !ok {
		r.Enum = requestStateDone
		return 0, nil
	}

	r.Body = append(r.Body, data...)

	numberRaw := length
	numberTrimmed := strings.TrimSpace(numberRaw)

	n, err := strconv.Atoi(numberTrimmed)
	if err != nil {
		return 0, errors.New("invalid content length")
	}

	r.BodyLength += len(data)

	if n == r.BodyLength {
		r.Enum = requestStateDone
	} else if r.BodyLength > n {
		return 0, errors.New("invalid content length")
	}
	return len(data), nil
}
