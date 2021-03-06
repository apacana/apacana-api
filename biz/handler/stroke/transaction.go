package stroke

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateUserStroke(c *gin.Context, userInfo *mysql.UserInfo, strokeList *transform.StrokeList, strokeName string) (strokeID int64, strokeToken string, createTime string, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "CreateUserStroke failed, err: %v", err)
			}
		}
	}()

	// create stroke
	createTime = time.Now().Format("2006-01-02 15:04:05")
	strokeToken = helper.GenerateToken([]byte{'s', 't', 'r', 'o', 'k', 'e'}, "")
	strokeInfo, err := mysql.InsertStrokeInfo(c, tx, &mysql.StrokeInfo{
		StrokeToken: strokeToken,
		StrokeName:  strokeName,
		OwnerID:     userInfo.ID,
		CreateTime:  createTime,
		UpdateTime:  createTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateUserStroke InsertStrokeInfo failed, err: %v", err)
		return
	}

	// update user
	if strokeList.DefaultStroke != 0 {
		strokeList.HistoryStrokeList = append(strokeList.HistoryStrokeList, strokeList.DefaultStroke)
	}
	strokeList.DefaultStroke = strokeInfo.ID
	strokeID = strokeInfo.ID
	err = mysql.UpdateUserInfo(c, tx, userInfo.ID, map[string]interface{}{
		"strokes": *transform.PackStrokeList(strokeList),
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateUserStroke UpdateUserInfo failed, err: %v, userID: %d", err, userInfo.ID)
		return
	}

	return
}
