package model

type LoginReq struct {
	EmpCode     string `json:"EmpCode" validate:"required"`
	EmpPassword string `json:"EmpPassword" validate:"required"`
}

type Emp struct {
	EmpName     string
	EmpCode     uint32
	EmpPassword string
	SecLevel    int32
	FixEmpStore int32
}

type LoginResponse struct {
	Token    string
	Employee Emp
}
