package model

type Order struct {
	DocNo     int
	DocDate   string
	EmpCode   int
	TotalCash float64
	EmpName   string
}
type OrderTotals struct {
	TotalCash     float64
	TotalPackages int
}

type InsertOrder struct {
	DocNo         int
	StoreCode     int
	EmpCode       int
	AccountSerial int
}

type GetOrderItemsRequest struct {
	Serial int `json:"Serial" validate:"required"`
}

type OrderItem struct {
	Serial   int
	BarCode  int
	ItemName string
	Qnt      float64
	Price    float64
	Total    float64
}

type InsertOrderItem struct {
	HeadSerial int
	ItemSerial int
	Qnt        float64
	Price      float64
}
