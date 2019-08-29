package user

import (
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/transform"
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

	touristStrokeList, err := transform.StringToStrokeList(touristInfo.Strokes)
	if err != nil {
		return
	}

	userStrokeList, err := transform.StringToStrokeList(userInfo.Strokes)
	if err != nil {
		return
	}

	// quota judge
	defaultNum := 0
	if touristStrokeList.DefaultStroke != 0 {
		defaultNum = 1
	}
	if len(userStrokeList.HistoryStrokeList)+len(touristStrokeList.HistoryStrokeList)+defaultNum > config.StrokeLimit {
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
	for _, touristStroke := range touristStrokeList.HistoryStrokeList {
		userStrokeList.HistoryStrokeList = append(userStrokeList.HistoryStrokeList, touristStroke)
	}
	if touristStrokeList.DefaultStroke != 0 {
		userStrokeList.HistoryStrokeList = append(userStrokeList.HistoryStrokeList, touristStrokeList.DefaultStroke)
	}

	strokeStr = *transform.PackStrokeList(userStrokeList)
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
