package food

import (
	"fmt"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func GetYelpFood(c *gin.Context) {
	yelpToken := c.Param("yelpToken")
	if yelpToken == "" {
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "GetYelpFood parm: %+v", yelpToken)

	url := fmt.Sprintf("https://api.yelp.com/v3/businesses/%s", yelpToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetYelpFood new request failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	req.Header.Set("Authorization", "Bearer bT-LePhROBnZoDkfr_ACjyneZzR51XvdUI59f3i5fzVahvjtquuRMcLGnuygCACEVoRY3oESxGh3C-rACx_yZ9aJq-9iYUL5g8PyZtARFV_rICNnQNNkSPJBnalJXXYx")
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetYelpFood http request failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetYelpFood http read body failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, string(body))
}

func SearchYelpFood(c *gin.Context) {
	var searchYelpFoodForm SearchYelpFoodForm
	if err := c.ShouldBindJSON(&searchYelpFoodForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "SearchYelpFood bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "SearchYelpFood from: %+v", searchYelpFoodForm)

	url := fmt.Sprintf("https://api.yelp.com/v3/businesses/search?term=food&latitude=%s&longitude=%s", searchYelpFoodForm.Latitude, searchYelpFoodForm.Longitude)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchYelpFood new request failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	req.Header.Set("Authorization", "Bearer bT-LePhROBnZoDkfr_ACjyneZzR51XvdUI59f3i5fzVahvjtquuRMcLGnuygCACEVoRY3oESxGh3C-rACx_yZ9aJq-9iYUL5g8PyZtARFV_rICNnQNNkSPJBnalJXXYx")
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchYelpFood http request failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchYelpFood http read body failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, string(body))
}
