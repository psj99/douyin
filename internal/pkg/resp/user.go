package resp

type UserRegisterResp struct {
	Response
	UserId int64  `json:"user_id"` // 用户id
	Token  string `json:"token"`   // 用户鉴权token
}

type UserLoginResp UserRegisterResp // 目前版本登录与注册响应元素类型完全相同

type UserInfoResp struct {
	Response
	User *UserInfo `json:"user"` // 用户信息
}

type UserInfo struct {
	UserId          int64  `json:"user_id"`           // 用户id
	Name            string `json:"name"`              // 用户名称
	FollowCount     int64  `json:"follow_count"`      // 关注总数
	FollowerCount   int64  `json:"follower_count"`    // 粉丝总数
	IsFollow        bool   `json:"is_follow" `        // true-已关注，false-未关注
	Avatar          string `json:"avatar"`            //用户头像
	BackGroundImage string `json:"background_image"`  //用户个人页顶部大图
	Signature       string `json:"signature "`        //个人简介
	TotalFavorited  string `json:"total_favorited  "` //获赞数量
	WorkCount       string `json:"work_count   "`     //作品数量
	FavoriteCount   string `json:"favorite_count   "` ///点赞数量
}
