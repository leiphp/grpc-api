package config

import (
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetConfig(t *testing.T) {
	c, err := GetConfig()
	assert.NoError(t, err)
	_, _ = pretty.Println(c)
}

func TestLoadTraceConfig(t *testing.T) {
	c := LoadTraceConfig()
	t.Logf("%+v\n", c)
}
