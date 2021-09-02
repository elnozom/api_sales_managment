package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.GET("/stores", h.GetStores)
	v1.GET("/account", h.GetAccount)
	v1.GET("/employee", h.GetEmp)
	//item

	v1.GET("/items", h.GetItems)
	v1.GET("/item", h.GetItem)
	v1.PUT("/item/update", h.UpdateItem)

	// order
	v1.GET("/orders", h.ListOrders)
	v1.POST("/orders", h.InsertOrder)
	v1.GET("/orders/no", h.GetSalesOrderDocNo)
	v1.POST("/orders/close", h.CloseOrder)
	v1.POST("/orders/item", h.InsertOrderItem)
	v1.PUT("/orders/item/update", h.UpdateOrderItem)
	v1.POST("/orders/item/delete", h.DeleteOrderItem)
	v1.GET("/orders/items", h.GetOrderItems)

}
