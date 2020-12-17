package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type people struct {
	Name string
}

func TestWriteToJSON(t *testing.T) {
	p := &people{Name: "jhx"}
	assert.Equal(t, "{\"Name\":\"jhx\"}", WriteToJSON(p))
}

func TestReadFromFile(t *testing.T) {
	str, err := ReadFromFile("test.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello test", str)
	str, err = ReadFromFile("123")
	assert.NotEqual(t, nil, err)
}
