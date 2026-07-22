package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "ea-qms"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error creating jwt token: %w", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			// validate the signing algorithm is what we expect
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
		jwt.WithIssuer(issuer),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}
	userIDAsString, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error retrieving User ID from verified token: %w", err)
	}
	userID, err := uuid.Parse(userIDAsString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing User ID from string to UUID type: %w", err)
	}
	return userID, nil
}
