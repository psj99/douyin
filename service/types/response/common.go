package response

type CommonResp struct {
	Status_Code int    `json:"status_code"` // 状态码，0-成功，其他值-失败
	Status_Msg  string `json:"status_msg"`  // 返回状态描述
}
