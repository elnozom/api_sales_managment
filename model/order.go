package model

type InsertOrder struct {
	DocNo     int
	StoreCode int
	EmpCode   int
}

type InsertOrderItem struct {
	HeadSerial int
	ItemSerial int
	Qnt        float64
	Price      float64
}
