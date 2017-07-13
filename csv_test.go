package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	r := bytes.NewReader([]byte("foo,bar,baz"))
	csv := NewCSVReader(r, ',')
	data, err := csv.Read()
	assert(t, err == nil, fmt.Sprintf("should not encounter error: %v", err))
	for i, el := range []string{"foo", "bar", "baz"} {
		assert(t, data[i] == el, fmt.Sprintf("element %d should be '%s': '%s'", i, el, data[i]))
	}
}

func TestReadWithDoubleQuotesSlashEscaped(t *testing.T) {
	r := bytes.NewReader([]byte(`test,"a \"b\" c"`))
	csv := NewCSVReader(r, ',')
	data, err := csv.Read()
	assert(t, err == nil, fmt.Sprintf("should not encounter error: %v", err))
	for i, el := range []string{"test", `a "b" c`} {
		assert(t, data[i] == el, fmt.Sprintf("element %d should be '%s': '%s'", i, el, data[i]))
	}
}

func TestReadWithDoubleQuotesDoubleEscaped(t *testing.T) {
	r := bytes.NewReader([]byte(`test,"a ""b"" c"`))
	csv := NewCSVReader(r, ',')
	data, err := csv.Read()
	assert(t, err == nil, fmt.Sprintf("should not encounter error: %v", err))
	for i, el := range []string{"test", `a "b" c`} {
		assert(t, data[i] == el, fmt.Sprintf("element %d should be '%s': '%s'", i, el, data[i]))
	}
}
