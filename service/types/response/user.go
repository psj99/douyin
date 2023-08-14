package response

type UserRegisterResp struct {
	Status
	User_Id uint   `json:"user_id"` // 用户id
	Token   string `json:"token"`   // 用户鉴权token
}

type UserLoginResp struct {
	Status
	User_Id uint   `json:"user_id"` // 用户id
	Token   string `json:"token"`   // 用户鉴权token
}

type UserInfoResp struct {
	Status
	User User `json:"user"` // 用户信息
}

// 用户信息/视频作者信息
type User struct {
	ID               uint   `json:"id"`               // 用户id
	Name             string `json:"name"`             // 用户名称
	Follow_Count     uint   `json:"follow_count"`     // 关注总数
	Follower_Count   uint   `json:"follower_count"`   // 粉丝总数
	Is_Follow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Avatar           string `json:"avatar"`           // 用户头像
	Background_Image string `json:"background_image"` // 用户个人页顶部大图
	Signature        string `json:"signature"`        // 个人简介
	Total_Favorited  string `json:"total_favorited"`  // 获赞数量
	Work_Count       uint   `json:"work_count"`       // 作品数
	Favorite_Count   uint   `json:"favorite_count"`   // 喜欢数
}
