package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const ContextUserIDKey = "UserID"

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUser(c *gin.Context) (UserID int64, err error) {
	uid, ok := c.Get("UserID")
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	UserID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return

}
func GetPageInfo(c *gin.Context) (int64, int64) {
	//获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
