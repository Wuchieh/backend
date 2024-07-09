package pkg

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
)

const (
	Google = "GUserInfo"
)

// GetGoogleUserinfo 取得 Userinfo
func GetGoogleUserinfo(c *gin.Context) (*oauth2.Userinfo, bool) {
	value, ok := c.Get(Google)
	if !ok {
		return nil, false
	}
	var info *oauth2.Userinfo
	info, ok = value.(*oauth2.Userinfo)
	if !ok {
		return nil, false
	}
	return info, true
}

// GetGoogleUserinfoUnsafe 取得 Userinfo
//
//	此方法不安全 建議使用 GetGoogleUserinfo
func GetGoogleUserinfoUnsafe(c *gin.Context) *oauth2.Userinfo {
	return c.Keys[Google].(*oauth2.Userinfo)
}
