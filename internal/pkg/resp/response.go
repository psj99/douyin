package resp

type Response struct {
	StatusCode int32  `json:"status_code"` // 状态码：0-成功, 其他-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

func ResponseOK() Response {
	return Response{
		StatusCode: 0,
		StatusMsg:  "成功",
	}
}

func ResponseErr(msg string) Response {
	return Response{
		StatusCode: 7,
		StatusMsg:  msg,
	}
}
