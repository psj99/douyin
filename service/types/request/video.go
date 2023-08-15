package request

type VideoFeedReq struct {
	Latest_Time string `json:"latest_time,omitempty" form:"latest_time"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token       string `json:"token,omitempty" form:"token"`             // 可选参数，用户登录状态下设置
}

type VideoPublishReq struct {
	//data FormFile 视频数据
	Token string `json:"token" form:"token" binding:"required"` // 用户鉴权token
	Title string `json:"title" form:"title" binding:"required"` // 视频标题
}

type VideoPublishListReq struct {
	Token   string `json:"token" form:"token" binding:"required"`     // 用户鉴权token
	User_ID string `json:"user_id" form:"user_id" binding:"required"` // 用户id
}
