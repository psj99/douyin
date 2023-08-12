package request

type UserRegisterReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"` // 注册用户名，最长32个字符
	Password string `json:"password" form:"password" binding:"required,max=32"` // 密码，最长32个字符
}

type UserLoginReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"` // 登录用户名
	Password string `json:"password" form:"password" binding:"required,max=32"` // 登录密码
}

type UserInfoReq struct {
	User_ID string `json:"user_id"` // 用户id
	Token   string `json:"token"`   // 用户鉴权token
}
