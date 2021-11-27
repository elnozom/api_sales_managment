package model

type Order struct {
	Serial         int
	StkTr01Serial  int
	DocNo          int
	DocDate        string
	EmpCode        int
	TotalCash      float64
	EmpName        string
	CustomerName   string
	CustomerCode   int
	CustomerSerial int
	Reserved       bool
}
type OrderTotals struct {
	TotalCash     float64
	TotalPackages float64
}

type InsertOrder struct {
	DocNo         int
	StoreCode     int
	AccountSerial int
}

type GetOrderItemsRequest struct {
	Serial int `json:"Serial" validate:"required"`
}

type OrderItem struct {
	Serial        int
	BarCode       int
	ItemName      string
	Qnt           float64
	Price         float64
	PriceMin      float64
	PriceMax      float64
	Total         float64
	QntAntherUnit float64
	AvgWeight     float64
}
type GetOrderItemsResonse struct {
	Items  []OrderItem
	Totals OrderTotals
}
type InsertOrderItem struct {
	HeadSerial    int
	ItemSerial    int
	Qnt           float64
	Price         float64
	QntAntherUnit float64
	PriceMax      float64
	PriceMin      float64
	MinorPerMajor int
}

type InsertOrderResp struct {
	Serial int
	No     int
}
