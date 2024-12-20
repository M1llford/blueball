package mysql

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"web2/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err := db.Select(&communityList, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Error("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityListByID 根据ID查询社区详情
func GetCommunityDetailByID(uid int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time,update_time
               from community
               where community_id = ?`
	if err := db.Get(communityDetail, sqlStr, uid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrInvalidId
			zap.L().Error("get community detail by id error", zap.Error(err))
		}
	}
	return communityDetail, err
}
