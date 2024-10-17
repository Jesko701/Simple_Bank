package token

import (
	"encoding/hex"

	"aidanwoods.dev/go-paseto"
)

func PasetoTokenMaker() (privateKeyHex string) {
	secretKey := paseto.NewV4AsymmetricSecretKey()

	privateKeyHex = hex.EncodeToString(secretKey.ExportBytes())

	return
}
