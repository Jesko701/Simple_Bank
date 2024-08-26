package token

type Maker interface {
	CreateToken(username string) (string, error)
	VerifyToken(token string) (*Payload, error)
}
