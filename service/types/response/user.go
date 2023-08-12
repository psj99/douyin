package response

type UserRegisterResp struct {
	Status_Code int    `json:"status_code"` // 状态码，0-成功，其他值-失败
	Status_Msg  string `json:"status_msg"`  // 返回状态描述
	User_Id     uint   `json:"user_id"`     // 用户id
	Token       string `json:"token"`       // 用户鉴权token
}

type UserLoginResp struct {
	Status_Code int    `json:"status_code"` // 状态码，0-成功，其他值-失败
	Status_Msg  string `json:"status_msg"`  // 返回状态描述
	User_Id     uint   `json:"user_id"`     // 用户id
	Token       string `json:"token"`       // 用户鉴权token
}
type UserInfoResp struct {
	Status_Code int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	Status_Msg  *string `json:"status_msg"`  // 返回状态描述
	User        *User   `json:"user"`        // 用户信息
}

// 用户信息/视频作者信息
type User struct {
	ID               int64  `json:"id"`               // 用户id
	Name             string `json:"name"`             // 用户名称
	Follow_Count     int64  `json:"follow_count"`     // 关注总数
	Follower_Count   int64  `json:"follower_count"`   // 粉丝总数
	Is_Follow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Avatar           string `json:"avatar"`           // 用户头像
	Background_Image string `json:"background_image"` // 用户个人页顶部大图
	Signature        string `json:"signature"`        // 个人简介
	Total_Favorited  string `json:"total_favorited"`  // 获赞数量
	Work_Count       int64  `json:"work_count"`       // 作品数
	Favorite_Count   int64  `json:"favorite_count"`   // 喜欢数
}
