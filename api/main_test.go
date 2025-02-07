package api

import (
	"os"
	db "simple-bank/db/sqlc"
	"simple-bank/utils"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		SecretKey: utils.RandomString(32),
		Duration:  time.Minute,
	}
	server, err := New(store, config)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
