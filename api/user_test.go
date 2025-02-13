package api

import (
	db "simple-bank/db/sqlc"
	"simple-bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) db.User {
	hash_pass, err := utils.CreateHashedPassword("secret")
	require.NoError(t, err)
	return db.User{
		Username:          utils.RandomOwner(),
		Email:             utils.RandomEmail(),
		HashedPassword:    hash_pass,
		CreatedAt:         time.Now(),
		FullName:          utils.RandomOwner(),
		PasswordChangedAt: time.Now(),
	}
}

func TestUser(t *testing.T) {

}
