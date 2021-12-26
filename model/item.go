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
	Serial             int
	Name               string
	Code               string
	AnQnt              float32
	ItemHaveAntherUnit bool
	Qnt                float32
	Price              float64
	LimitedQnt         float64
	StopSale           bool
	PMin               float64
	PMax               float64
	AvrWeight          float64
	MinorPerMajor      int
}

type IetmBalance struct {
	Raseed           float64
	ItemName         string
	RaseedReserved   float64
	RaseedNet        float64
	AnRaseed         float64
	AnRaseedReserved float64
	AnRaseedNet      float64
	StoreCode        float64
	StoreName        string
	DsiplayName      string
}
