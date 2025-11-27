package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HOST: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitspace
	headers = NewHeaders()
	data = []byte("Host:    localhost:42069        \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 34, n)
	assert.False(t, done)

	// Test Multiple Values

	headers = NewHeaders()
	data1 := []byte("Content-Type: application/json\r\n\r\n")
	data2 := []byte("Authorization: Bearer token123\r\n\r\n")
	data3 := []byte("Host: localhost:42069\r\n\r\n")
	_, _, _ = headers.Parse(data1)
	_, _, _ = headers.Parse(data2)
	n, done, err = headers.Parse(data3)

	require.NoError(t, err)
	assert.False(t, done)
	assert.Greater(t, n, 0)

	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, "Bearer token123", headers["authorization"])
	assert.Equal(t, "localhost:42069", headers["host"])

}
