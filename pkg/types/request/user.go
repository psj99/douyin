package request

type UserRegisterReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"`
	Password string `json:"password" form:"password" binding:"required,max=32"`
}

type UserLoginReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"`
	Password string `json:"password" form:"password" binding:"required,max=32"`
}
