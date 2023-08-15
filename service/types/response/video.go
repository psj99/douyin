package response

type VideoFeedResp struct {
	Status
	Next_Time  uint    `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	Video_List []Video `json:"video_list"` // 视频列表
}

type VideoPublishResp struct {
	Status
}

type VideoPublishListResp struct {
	Status
	Video_List []Video `json:"video_list"` // 用户发布的视频列表
}

// 视频列表
type Video struct {
	ID             uint   `json:"id"`             // 视频唯一标识
	Author         User   `json:"author"`         // 视频作者信息
	Play_URL       string `json:"play_url"`       // 视频播放地址
	Cover_URL      string `json:"cover_url"`      // 视频封面地址
	Favorite_Count uint   `json:"favorite_count"` // 视频的点赞总数
	Comment_Count  uint   `json:"comment_count"`  // 视频的评论总数
	Is_Favorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title          string `json:"title"`          // 视频标题
}
