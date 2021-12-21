package handler

import (
	"hand_held/config"
	"hand_held/router/middleware"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	jwtMiddleware := middleware.JWT(config.Config("JWT_SECRET"))
	auth := v1.Group("/", jwtMiddleware)
	auth.GET("stores", h.GetStores)
	auth.GET("account", h.GetAccount)
	v1.POST("/login", h.Login)
	auth.GET("employee", h.GetEmp)
	//item
	auth.GET("items", h.GetItems)
	auth.GET("item", h.GetItem)
	auth.GET("item/balnace/:serial", h.GetItemBalance)
	auth.PUT("item/update", h.UpdateItem)

	// order
	auth.PUT("unreserve", h.UpdateReservedForEmp, jwtMiddleware)
	auth.GET("orders", h.ListOrders)
	auth.GET("orders/store", h.ListStoreOrders, jwtMiddleware)
	auth.GET("orders/store/:Serial", h.ListStoreOrderItems, jwtMiddleware)
	auth.PUT("orders/store/update", h.UpdateStoreOrderItems, jwtMiddleware)
	auth.PUT("orders/store/close/:serial", h.CloseStoreOrderItems, jwtMiddleware)

	auth.PUT("orders/update/:Serial", h.UpdateOrder)
	auth.POST("orders", h.InsertOrder)
	auth.GET("orders/no", h.GetSalesOrderDocNo)
	auth.POST("orders/close", h.CloseOrder)
	auth.POST("orders/item", h.InsertOrderItem)
	auth.PUT("orders/item/update", h.UpdateOrderItem)
	auth.POST("orders/item/delete", h.DeleteOrderItem)
	auth.GET("orders/items", h.GetOrderItems)

}
