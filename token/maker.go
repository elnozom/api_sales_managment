package token

import "time"

// maker is an interface for manging tokens
type Maker interface {
	// create a new token for sepecifi employee
	CreateToken(empCode uint32, duration time.Duration) (string, error)

	// verifuToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
