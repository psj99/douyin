package response

type VideoFeedResp struct {
	Status_Code int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	Status_Msg  *string `json:"status_msg"`  // 返回状态描述
	Next_Time   *int64  `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	Video_List  []Video `json:"video_list"`  // 视频列表
}

// 视频列表
type Video struct {
	ID             int64  `json:"id"`             // 视频唯一标识
	Author         User   `json:"author"`         // 视频作者信息
	Play_URL       string `json:"play_url"`       // 视频播放地址
	Cover_URL      string `json:"cover_url"`      // 视频封面地址
	Favorite_Count int64  `json:"favorite_count"` // 视频的点赞总数
	Comment_Count  int64  `json:"comment_count"`  // 视频的评论总数
	Is_Favorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title          string `json:"title"`          // 视频标题
}
