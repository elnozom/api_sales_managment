package model

type OrderResp struct {
	Serial    int
	DocNo     string
	DocDate   string
	Discount  float64
	TotalCash float64
	TotalTax  float64
}
