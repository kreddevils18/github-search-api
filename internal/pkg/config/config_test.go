package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig()
	require.Nil(t, err)
	require.NotNil(t, conf.Github)
	require.NotNil(t, conf.DB)
}
