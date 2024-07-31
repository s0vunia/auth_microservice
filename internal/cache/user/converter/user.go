package converter

import (
	"database/sql"
	"time"

	modelCache "github.com/s0vunia/auth_microservices_course_boilerplate/internal/cache/user/model"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

func ToUserCacheFromModel(user *model.User) *modelCache.User {
	var updatedAtNs *int64
	if user.UpdatedAt.Valid {
		updatedAtNs = new(int64)
		*updatedAtNs = user.UpdatedAt.Time.Unix()
	}
	return &modelCache.User{
		ID:          user.ID,
		Name:        user.Info.Name,
		Email:       user.Info.Email,
		Role:        modelCache.Role(user.Info.Role),
		CreatedAtNs: user.UpdatedAt.Time.Unix(),
		UpdatedAtNs: updatedAtNs,
	}
}

func ToUserFromCache(user *modelCache.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromCache(user),
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromCache(info *modelCache.User) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  model.Role(info.Role),
	}
}
