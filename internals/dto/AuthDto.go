package dto

type LoginReq struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LogInRes struct {
	UserId string `json:"userId,omitempty"`
	Token  string `json:"token,omitempty"`
}
