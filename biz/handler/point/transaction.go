package point

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/out"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"time"
)

func addStrokePointList(c *gin.Context, strokeInfo *mysql.StrokeInfo, pointList *transform.PointList, form AddPointForm, pointType mysql.PointType) (outPut map[string]interface{}, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "addStrokePointList failed, err: %v", err)
			}
		}
	}()

	// create point
	iconType := ""
	if form.IconType != nil && len(*form.IconType) < 50 {
		iconType = *form.IconType
	}
	iconColor := ""
	if form.IconColor != nil && len(*form.IconColor) < 50 {
		iconColor = *form.IconColor
	}
	ext := ""
	if form.Ext != nil {
		ext = *form.Ext
	}
	placeName := ""
	if form.PlaceName != nil && len(*form.PlaceName) < 301 {
		placeName = *form.PlaceName
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	pointToken := helper.GenerateToken([]byte{'p', 'o', 'i', 'n', 't'}, "")
	err = mysql.InsertPointInfo(c, tx, &mysql.PointInfo{
		PointToken: pointToken,
		PointID:    form.PointID,
		PointType:  pointType,
		Text:       form.Text,
		PlaceName:  placeName,
		Center:     form.Center,
		Comment:    "",
		IconType:   iconType,
		IconColor:  iconColor,
		Ext:        ext,
		StrokeID:   strokeInfo.ID,
		CreateTime: nowTime,
		UpdateTime: nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "addStrokePointList InsertPointInfo failed, err: %v", err)
		return
	}
	pointInfo, err := mysql.GetPointByToken(c, tx, pointToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "addStrokePointList GetPointByToken failed, err: %v, pointToken: %v", err, pointToken)
		return
	}
	pointInfoOut := &out.PointInfoOut{
		PointToken: pointInfo.PointToken,
		PointID:    pointInfo.PointID,
		PointType:  form.PointType,
		Text:       pointInfo.Text,
		PlaceName:  pointInfo.PlaceName,
		Center:     pointInfo.Center,
		Comment:    pointInfo.Comment,
		IconType:   pointInfo.IconType,
		IconColor:  pointInfo.IconColor,
		Ext:        pointInfo.Ext,
	}

	// update stroke
	pointList.PointList = append(pointList.PointList, pointInfo.ID)
	err = mysql.UpdateStrokeByToken(c, tx, strokeInfo.StrokeToken, map[string]interface{}{
		"points_list": *transform.PackPointList(pointList),
		"update_time": nowTime,
	})

	outPut = map[string]interface{}{
		"stroke_info": out.StrokeUpdateOut{
			StrokeToken: strokeInfo.StrokeToken,
			StrokeName:  strokeInfo.StrokeName,
			UpdateTime:  nowTime,
		},
		"point_info": pointInfoOut,
	}

	return
}

func deleteStrokePointList(c *gin.Context, strokeInfo *mysql.StrokeInfo, pointList *transform.PointList, pointInfo *mysql.PointInfo) (outPut map[string]interface{}, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "deleteStrokePointList failed, err: %v", err)
			}
		}
	}()

	// set point delete
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = mysql.UpdatePointByID(c, tx, pointInfo.ID, map[string]interface{}{
		"status":      helper.PointDeleteStatus,
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "deleteStrokePointList UpdatePointByID failed, err: %v", err)
		return
	}

	// update stroke
	pointList.PointList = helper.ArrayRemove(pointList.PointList, pointInfo.ID)
	err = mysql.UpdateStrokeByToken(c, tx, strokeInfo.StrokeToken, map[string]interface{}{
		"points_list": *transform.PackPointList(pointList),
		"update_time": nowTime,
	})

	outPut = map[string]interface{}{
		"stroke_info": out.StrokeUpdateOut{
			StrokeToken: strokeInfo.StrokeToken,
			StrokeName:  strokeInfo.StrokeName,
			UpdateTime:  nowTime,
		},
	}

	return
}

func recreateStrokePointList(c *gin.Context, strokeInfo *mysql.StrokeInfo, pointList *transform.PointList, pointInfo *mysql.PointInfo) (outPut map[string]interface{}, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "recreateStrokePointList failed, err: %v", err)
			}
		}
	}()

	// update point
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = mysql.UpdatePointByID(c, tx, pointInfo.ID, map[string]interface{}{
		"status":      helper.PointNormalStatus,
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "recreateStrokePointList UpdatePointByID failed, err: %v", err)
		return
	}
	pointTypeName, _ := helper.GetNameByPointType(pointInfo.PointType)
	pointInfoOut := &out.PointInfoOut{
		PointToken: pointInfo.PointToken,
		PointID:    pointInfo.PointID,
		PointType:  pointTypeName,
		Text:       pointInfo.Text,
		PlaceName:  pointInfo.PlaceName,
		Center:     pointInfo.Center,
		Comment:    pointInfo.Comment,
		IconType:   pointInfo.IconType,
		IconColor:  pointInfo.IconColor,
		Ext:        pointInfo.Ext,
	}

	// update stroke
	pointList.PointList = append(pointList.PointList, pointInfo.ID)
	err = mysql.UpdateStrokeByToken(c, tx, strokeInfo.StrokeToken, map[string]interface{}{
		"points_list": *transform.PackPointList(pointList),
		"update_time": nowTime,
	})

	outPut = map[string]interface{}{
		"stroke_info": out.StrokeUpdateOut{
			StrokeToken: strokeInfo.StrokeToken,
			StrokeName:  strokeInfo.StrokeName,
			UpdateTime:  nowTime,
		},
		"point_info": pointInfoOut,
	}

	return
}
