package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

func CreateToken(payload jwt.MapClaims, secret string) (string, error) {
	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetJwtToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CreatePayload(userID int, username, role string, exp int64, iss string) jwt.MapClaims {
	payload := jwt.MapClaims{}
	payload["user_id"] = userID
	payload["username"] = username
	payload["role"] = role
	payload["exp"] = exp
	payload["iss"] = iss
	return payload
}

func ValidateJwtToken(token, secret string) (jwt.MapClaims, error) {
	jwtToken, err := GetJwtToken(token, secret)
	if err != nil {
		return nil, fmt.Errorf("err get token: %w", err)
	}
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error getting claims")
	}
	return claims, nil
}
