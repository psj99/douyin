package response

type UserRegisterResp struct {
	Status_Code int    `json:"status_code"`
	Status_Msg  string `json:"status_msg"`
	User_Id     uint   `json:"user_id"`
	Token       string `json:"token"`
}

type UserLoginResp UserRegisterResp // 目前版本登录与注册响应元素类型完全相同
