package response

type FeedResp struct {
	StatusCode int32    `json:"status_code" binding:"required"` // 状态码：0-成功, 其他-失败
	StatusMsg  string   `json:"status_msg"`                     // 返回状态描述
	VideoList  []*Video `json:"video_list"`                     // 视频列表
	NextTime   int64    `json:"next_time"`                      // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

type PublishActionResp struct {
	StatusCode int32  `json:"status_code" binding:"required"` // 状态码：0-成功, 其他-失败
	StatusMsg  string `json:"status_msg"`                     // 返回状态描述
}

type PublishListResp struct {
	StatusCode int32    `json:"status_code"  binding:"required"` // 状态码，0-成功，其他值-失败
	StatusMsg  string   `json:"status_msg"`                      // 状态码，0-成功，其他值-失败
	VideoList  []*Video `json:"video_list"`                      // 用户发布的视频列表
}

type Video struct {
	Id            int64     // 视频唯一标识
	User          *UserInfo // 视频作者信息
	PlayUrl       string    // 视频播放地址
	CoverUrl      string    // 视频封面地址
	FavoriteCount int64     // 视频的点赞总数
	CommentCount  int64     // 视频的评论总数
	IsFavorite    bool      // true-已点赞，false-未点赞
	Title         string    // 视频标题
}
