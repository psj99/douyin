package response

type UserRegisterResp struct {
	Code   int    `json:"status_code"`
	Msg    string `json:"status_msg"`
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserLoginResp struct {
	Code   int    `json:"status_code"`
	Msg    string `json:"status_msg"`
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}
