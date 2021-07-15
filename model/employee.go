package model

type EmpReq struct {
	EmpCode int `json:"EmpCode" validate:"required"`
}

type Emp struct {
	EmpName     string
	EmpCode     int
	EmpPassword string
}
