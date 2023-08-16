package dao

import (
	"context"
)

// 检查是否关注
func CheckFollow(ctx context.Context, user_id uint, follow_id uint) (is_follower bool, err error) {
	if user_id == follow_id {
		return false, nil // 默认自己不关注自己
	}

	DB := GetDB(ctx)
	var count int64
	err = DB.Table("follow").Where("user_id=? AND follow_id=?", user_id, follow_id).Count(&count).Error
	if err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}
