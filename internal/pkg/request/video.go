package request

type PublishActionReq struct {
	Token string `binding:"required"` // 用户鉴权token
	Data  []byte `binding:"required"` // 视频数据
	Title string `binding:"required"` // 视频标题
}

type FeedReq struct {
	LatestTime int64  // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时
	Token      string // 可选参数，登录用户设置
}
