package token

import "time"

//maker interface is for managing tokens

type Maker interface {
	//CreateToken creates a new token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken verifies the generated token
	VerifyToken(token string) (*Payload, error)
}
