package model

type LoginReq struct {
	EmpCode     uint32 `json:"EmpCode" validate:"required"`
	EmpPassword string `json:"EmpPassword" validate:"required"`
}

type Emp struct {
	EmpName     string
	EmpCode     uint32
	EmpPassword string
}

type LoginResponse struct {
	token    string
	employee Emp
}
