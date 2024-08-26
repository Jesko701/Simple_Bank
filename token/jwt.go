package token

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	privateKey *ecdsa.PrivateKey //ECDSA private key
	publicKey  *ecdsa.PublicKey  //ECDSA private key
}

func NewJWTMaker(privateKeyString string) (Maker, error) {
	// Decode the PEM block
	block, _ := pem.Decode([]byte(privateKeyString))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key: %v", privateKeyString)
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Assert the private key type
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type *ecdsa.PrivateKey")
	}

	return &JWTMaker{
		privateKey: ecdsaPrivateKey,
		publicKey:  &ecdsaPrivateKey.PublicKey,
	}, nil
}

func (j *JWTMaker) CreateToken(username string) (string, error) {
	payload, err := NewPayload(username)
	if err != nil {
		return "", err
	}

	// Create a new token with the payload claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":       payload.ID,
		"username": payload.Username,
		"iat":      payload.IssuedAt.Unix(),
		"exp":      payload.ExpiredAt.Unix(),
	})

	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method conforms to "SigningMethodECDSA"
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpecte sign-in method")
		}
		return j.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert Unix timestamps from claims to time.Time
		issuedAt, err := getUnixTimeStamps(claims["iat"])
		if err != nil {
			return nil, fmt.Errorf("invalid type for claim: iat is invalid")
		}
		expiredAt, err := getUnixTimeStamps(claims["exp"])
		if err != nil {
			return nil, fmt.Errorf("invalid type for claim: exp is invalid")
		}
		return &Payload{
			ID:        uuid.Must(uuid.Parse(claims["id"].(string))),
			Username:  claims["username"].(string),
			IssuedAt:  issuedAt,
			ExpiredAt: expiredAt,
		}, nil
	}

	return nil, ErrInvalidToken
}

func getUnixTimeStamps(claims interface{}) (time.Time, error) {
	timestamps, ok := claims.(float64)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid type for claim")
	}
	return time.Unix(int64(timestamps), 0), nil
}
