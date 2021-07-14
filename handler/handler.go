package handler

	
import "github.com/jinzhu/gorm"

type Handler struct {
	db *gorm.DB
}

func NewHandler(databaase *gorm.DB) *Handler {
	return &Handler{
		db : databaase,
	}
}
