package handler

import (
	"hand_held/config"
	"hand_held/router/middleware"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	jwtMiddleware := middleware.JWT(config.Config("JWT_SECRET"))
	v1.GET("/stores", h.GetStores)
	v1.GET("/account", h.GetAccount)

	v1.POST("/login", h.Login)
	auth := v1.Group("/", jwtMiddleware)
	auth.GET("employee", h.GetEmp)
	//item
	v1.GET("/items", h.GetItems)
	v1.GET("/item", h.GetItem)
	v1.PUT("/item/update", h.UpdateItem)

	// order
	auth.GET("orders", h.ListOrders)
	auth.POST("orders/exit", h.ExitOrder)
	auth.POST("orders", h.InsertOrder)
	v1.GET("/orders/no", h.GetSalesOrderDocNo)
	v1.POST("/orders/close", h.CloseOrder)
	v1.POST("/orders/item", h.InsertOrderItem)
	v1.PUT("/orders/item/update", h.UpdateOrderItem)
	v1.POST("/orders/item/delete", h.DeleteOrderItem)
	v1.GET("/orders/items", h.GetOrderItems)

}
