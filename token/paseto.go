package token

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	privateKey paseto.V4AsymmetricSecretKey
	publicKey  paseto.V4AsymmetricPublicKey
}

func NewPasetoMaker(PrivateKey string) (Maker, error) {
	// Decode the private key from hex
	privateKeyBytes, err := hex.DecodeString(PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()

	return &PasetoMaker{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// Parameter duration is not used in paseto
func (m *PasetoMaker) CreateToken(username string) (string, error) {
	payload, err := NewPayload(username)
	if err != nil {
		return "", err
	}

	token := paseto.NewToken()
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)
	token.SetString("id", payload.ID.String())
	token.SetString("username", payload.Username)

	signedToken := token.V4Sign(m.privateKey, nil)
	log.Printf("Created Verifying Token: %v\n", token)

	return signedToken, nil
}

func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload

	parsedToken, err := paseto.NewParser().ParseV4Public(m.publicKey, token, nil)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	idStr, err := parsedToken.GetString("id")
	if err != nil {
		fmt.Println("Error getting ID from token:", err)
		return nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("Error parsing UUID from token:", err)
		return nil, err
	}

	username, err := parsedToken.GetString("username")
	if err != nil {
		fmt.Println("Error getting username from token:", err)
		return nil, err
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		fmt.Println("Error getting issued at time from token:", err)
		return nil, err
	}

	expiration, err := parsedToken.GetExpiration()
	if err != nil {
		fmt.Println("Error getting expiration time from token:", err)
		return nil, err
	}

	log.Printf("Token expiration: %v", expiration)

	payload = Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiration,
	}

	log.Printf("Token issued at: %v", issuedAt)
	log.Printf("Token expiration: %v", expiration)
	log.Printf("Current time: %v", time.Now())

	if err := payload.Valid(); err != nil {
		fmt.Println("Payload validation error:", err)
		return nil, err
	}

	return &payload, nil
}
