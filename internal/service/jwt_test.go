package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var claimExample = jwt.MapClaims{
	"user_id":    "e4b67dfd-e55f-4c68-8742-08f0750c038e",
	"username":   "+77017284345",
	"role":       "owner",
	"type":       "access",
	"company_id": "afb6d5c9-5f53-46a6-80bb-e6ccea227111",
	"exp":        100,
	"iss":        "test",
}

func TestValidateJwtToken(t *testing.T) {
	claim := claimExample
	claimExample["exp"] = time.Now().Add(time.Second * time.Duration(1800)).Unix()
	secret := "xrt57g2jj8mj"

	t.Run("Success", func(t *testing.T) {
		accessToken, err := CreateToken(claim, secret)
		require.NoError(t, err)

		token, err := GetJwtToken(accessToken, secret)
		require.NoError(t, err)

		claim, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)
		resultClaim, err := ValidateJwtToken(accessToken, secret)
		require.NoError(t, err)

		assert.Equal(t, resultClaim, claim)
	})

	t.Run("Wrong secret", func(t *testing.T) {
		accessToken, err := CreateToken(claim, secret)
		require.NoError(t, err)

		resultClaim, err := ValidateJwtToken(accessToken, "wrong_secret")
		require.Error(t, err)
		assert.Nil(t, resultClaim)
	})
}
