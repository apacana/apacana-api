package user

import (
	"errors"
	"fmt"
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
)

func userStrokeTrans(c *gin.Context, touristInfo *mysql.UserInfo, userInfo *mysql.UserInfo) (strokeStr string, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "createUserStroke failed, err: %v", err)
			}
		}
	}()

	touristStrokeList := &helper.StrokeList{
		StrokeList: make([]int64, 0),
	}
	if len(touristInfo.Strokes) != 0 {
		touristStrokeList = helper.StringToStrokeList(touristInfo.Strokes)
		if touristStrokeList == nil {
			err = errors.New(fmt.Sprintf("StringToStrokeList error, strokes: %v", touristInfo.Strokes))
			return
		}
	}
	userStrokeList := &helper.StrokeList{
		StrokeList: make([]int64, 0),
	}
	if len(userInfo.Strokes) != 0 {
		userStrokeList = helper.StringToStrokeList(userInfo.Strokes)
		if userStrokeList == nil {
			err = errors.New(fmt.Sprintf("StringToStrokeList error, strokes: %v", userInfo.Strokes))
			return
		}
	}

	// quota judge
	if len(userStrokeList.StrokeList)+len(touristStrokeList.StrokeList) > config.StrokeLimit {
		err = helper.ErrStrokeOutOfLimit
		return
	}

	// change ownerID
	err = mysql.ChangeStrokeOwner(c, tx, touristInfo.ID, userInfo.ID)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "userStrokeTrans ChangeStrokeOwner failed, err: %v, touristID: %d, userID: %d", err, touristInfo.ID, userInfo.ID)
		return
	}

	// update user
	for _, touristStroke := range touristStrokeList.StrokeList {
		userStrokeList.StrokeList = append(userStrokeList.StrokeList, touristStroke)
	}

	strokeStr = *helper.PackStrokeList(userStrokeList)
	err = mysql.UpdateUserInfo(c, tx, userInfo.ID, map[string]interface{}{
		"strokes": strokeStr,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "userStrokeTrans UpdateUserInfo failed, err: %v, userID: %d", err, userInfo.ID)
		return
	}
	err = mysql.UpdateUserInfo(c, tx, touristInfo.ID, map[string]interface{}{
		"strokes": "",
		"status":  helper.TransferredStatus,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "userStrokeTrans UpdateUserInfo failed, err: %v, userID: %d", err, userInfo.ID)
		return
	}
	return
}
