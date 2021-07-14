package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.GET("/stores", h.GetStores)
	v1.GET("/account", h.GetAccount)
	v1.GET("/items", h.GetItems)
	v1.GET("/item", h.GetItem) // done

}
