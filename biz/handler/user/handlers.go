package user

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func CreateTourist(c *gin.Context) {
	var createTouristForm CreateTouristForm
	if err := c.ShouldBindJSON(&createTouristForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "CreateTourist bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "CreateTourist from: %+v", createTouristForm)

	token := helper.GenerateToken([]byte{'u', 's', 'e', 'r'}, createTouristForm.Key)
	cookie := helper.SetCookie(token, helper.SessionSalt)

	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		helper.ApacanaSession: cookie,
	})
}

func GetUserInfo(c *gin.Context) {
	userToken := c.GetString(helper.UserToken)
	helper.FormatLogPrint(helper.LOG, "GetUserInfo userToken: %+v", userToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
				"is_tourist": true,
			})
			return
		}
		helper.FormatLogPrint(helper.ERROR, "GetUserInfo GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	strokeInfoList, err := helper.CreateFmtStrokeList(c, userInfo.Strokes)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetUserInfo CreateFmtStrokeList failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	if userInfo.Status == 0 {
		helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
			"is_tourist":  true,
			"stroke_list": strokeInfoList,
		})
		return
	}
	newSession := helper.SetCookie(userToken, helper.SessionSalt)
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		"is_tourist": false,
		"user_info": map[string]interface{}{
			"name":                userInfo.Name,
			"token":               userInfo.Token,
			"status":              userInfo.Status,
			helper.ApacanaSession: newSession,
		},
		"stroke_list": strokeInfoList,
	})
}

func RegisterUser(c *gin.Context) {
	var registerUserForm RegisterUserForm
	if err := c.ShouldBindJSON(&registerUserForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "RegisterUser bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	if len(registerUserForm.PassWord) < 8 || len(registerUserForm.UserName) < 4 ||
		len(registerUserForm.UserName) > 16 || len(registerUserForm.Name) == 0 {
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "RegisterUser from: %+v", registerUserForm)

	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// insert
			nowTime := time.Now().Format("2006-01-02 15:04:05")
			passWord := helper.Md5(registerUserForm.PassWord)
			err := mysql.InsertUserInfo(c, &mysql.UserInfo{
				Token:      userToken,
				UserName:   registerUserForm.UserName,
				PassWord:   passWord,
				Name:       registerUserForm.Name,
				Status:     helper.LoginUserStatus,
				CreateTime: nowTime,
				UpdateTime: nowTime,
			})
			if err != nil {
				helper.FormatLogPrint(helper.ERROR, "RegisterUser InsertUserInfo failed, err: %v, userToken: %v", err, userToken)
				helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
				return
			}
			newSession := helper.SetCookie(userToken, helper.SessionSalt)
			helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
				"user_info": map[string]interface{}{
					"name":                registerUserForm.Name,
					"token":               userToken,
					"status":              helper.LoginUserStatus,
					helper.ApacanaSession: newSession,
				},
			})
			return
		}
		helper.FormatLogPrint(helper.ERROR, "RegisterUser GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if userInfo.Status != helper.TouristStatus {
		helper.FormatLogPrint(helper.WARNING, "RegisterUser not tourist userToken: %v", userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	// update
	passWord := helper.Md5(registerUserForm.PassWord)
	attrs := map[string]interface{}{
		"status":    helper.LoginUserStatus,
		"user_name": registerUserForm.UserName,
		"pass_word": passWord,
		"name":      registerUserForm.Name,
	}
	err = mysql.UpdateUserInfo(c, nil, userInfo.ID, attrs)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "RegisterUser UpdateUserInfo failed, err: %v, userToken: %v", err, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	newSession := helper.SetCookie(userToken, helper.SessionSalt)
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		"user_info": map[string]interface{}{
			"name":                registerUserForm.Name,
			"token":               userInfo.Token,
			"status":              helper.LoginUserStatus,
			helper.ApacanaSession: newSession,
		},
	})
}

func LoginUser(c *gin.Context) {
	var loginUserForm LoginUserForm
	if err := c.ShouldBindJSON(&loginUserForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "LoginUser bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "LoginUser from: %+v", loginUserForm)
	passWord := helper.Md5(loginUserForm.PassWord)
	userInfo, err := mysql.GetUserByUserPassWord(c, nil, loginUserForm.UserName, passWord)
	if err != nil && err != gorm.ErrRecordNotFound {
		helper.FormatLogPrint(helper.ERROR, "LoginUser GetUserByUserPassWord failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if userInfo == nil || userInfo.Status == helper.TouristStatus {
		helper.FormatLogPrint(helper.WARNING, "LoginUser status not 0")
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	strokeStr := userInfo.Strokes

	touristToken := c.GetString(helper.UserToken)
	if userInfo.Token == touristToken {
		strokeInfoList, err := helper.CreateFmtStrokeList(c, strokeStr)
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "LoginUser CreateFmtStrokeList failed, err: %v", err)
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		}
		newSession := helper.SetCookie(userInfo.Token, helper.SessionSalt)
		helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
			"user_info": map[string]interface{}{
				"name":                userInfo.Name,
				"token":               userInfo.Token,
				"status":              helper.LoginUserStatus,
				helper.ApacanaSession: newSession,
				"stroke_list":         strokeInfoList,
			},
		})
		return
	}

	touristInfo, err := mysql.GetUserInfoByToken(c, nil, touristToken)
	if err != nil && err != gorm.ErrRecordNotFound {
		helper.FormatLogPrint(helper.ERROR, "LoginUser GetUserInfoByToken failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	if touristInfo != nil {
		newStrokeStr, err := userStrokeTrans(c, touristInfo, userInfo)
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "LoginUser userStrokeTrans failed, err: %v", err)
			if err == helper.ErrStrokeOutOfLimit {
				helper.BizResponse(c, http.StatusOK, helper.CodeStrokeOutOfLimit, nil)
				return
			}
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		}
		strokeStr = newStrokeStr
	}
	strokeInfoList, err := helper.CreateFmtStrokeList(c, strokeStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "LoginUser CreateFmtStrokeList failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	newSession := helper.SetCookie(userInfo.Token, helper.SessionSalt)
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		"user_info": map[string]interface{}{
			"name":                userInfo.Name,
			"token":               userInfo.Token,
			"status":              helper.LoginUserStatus,
			helper.ApacanaSession: newSession,
		},
		"stroke_list": strokeInfoList,
	})
}
