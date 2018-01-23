package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader(`
development:
  datasource: root@localhost/dev

test:
  datasource: root@localhost/test
`)

	configs, err := NewConfigs(r)
	assert.NoError(err)
	c, ok := configs["development"]
	assert.True(ok)
	assert.Equal("root@localhost/dev", c.DSN(), "they should be equal")
}
