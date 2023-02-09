package dbcontext

import (
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMongo(t *testing.T) {
	conf, _ := config.LoadConfig(config.GetPath())
	_, err := NewMongo(conf)
	require.Nil(t, err)
}
