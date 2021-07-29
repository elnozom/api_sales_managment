package model

type GetItemRequest struct {
	BCode     string `json:"BCode" validate:"required"`
	Name      string `json:"Name"`
	StoreCode int    `json:"StoreCode" validate:"required"`
}

type SingleItem struct {
	Serial            int
	ItemName          string
	MinorPerMajor     int
	POSPP             float64
	POSTP             float64
	ByWeight          bool
	WithExp           bool
	ItemHasAntherUnit bool
}

type Item struct {
	Serial     int
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
