package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload struct will contains the payload data of the token

type Payload struct {
	ID        uuid.UUID `json:"id"`
	EmpCode   uint32    `json:"empCode"`
	IssudedAt time.Time `json:"essudedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

// will  create a new token payload to specific
func NewPayload(empCode uint32, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenId,
		EmpCode:   empCode,
		IssudedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

//Valid checks if the token is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
