package model

type GetItemRequest struct {
	BCode     string `json:"BCode" validate:"required"`
	StoreCode int    `json:"StoreCode" validate:"required"`
}

type SingleItem struct {
	Serial        int
	ItemName      string
	MinorPerMajor int
	POSPP         float64
	POSTP         float64
	ByWeight      bool
}

type Item struct {
	Name       string
	Code       string
	AnQnt      float32
	Qnt        float32
	Price      float64
	LimitedQnt bool
	StopSale   bool
	PMin       float64
	PMax       float64
}
