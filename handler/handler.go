package handler

import (
	"fmt"
	"hand_held/config"
	"hand_held/token"

	"github.com/jinzhu/gorm"
)

type Handler struct {
	db         *gorm.DB
	tokenMaker token.Maker
}

func NewHandler(databaase *gorm.DB) (*Handler, error) {
	tokenMaker, err := token.NewJWTMaker(config.Config("JWT_SECRET"))
	if err != nil {
		return nil, fmt.Errorf("cant create token maker: %w", err)
	}
	return &Handler{
		db:         databaase,
		tokenMaker: tokenMaker,
	}, nil
}
