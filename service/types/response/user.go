package response

type UserRegisterResp struct {
	Status_Code int    `json:"status_code"`
	Status_Msg  string `json:"status_msg"`
	User_Id     uint   `json:"user_id"`
	Token       string `json:"token"`
}

type UserLoginResp struct {
	Status_Code int    `json:"status_code"`
	Status_Msg  string `json:"status_msg"`
	User_Id     uint   `json:"user_id"`
	Token       string `json:"token"`
}
