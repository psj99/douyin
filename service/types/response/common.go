package response

type CommonResp struct {
	Status_Code int    `json:"status_code"`
	Status_Msg  string `json:"status_msg"`
}
