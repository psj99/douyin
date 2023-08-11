package request

type UserRegisterReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"`
	Password string `json:"password" form:"password" binding:"required,max=32"`
}

type UserLoginReq UserRegisterReq // 目前版本登录与注册请求元素类型完全相同
