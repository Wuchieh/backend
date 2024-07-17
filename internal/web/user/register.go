package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"string_backend_0001/internal/logger"
	"string_backend_0001/internal/model"
	"string_backend_0001/pkg"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @summary user register
// @description user register
// @tags User
// @id register
// @accept json
// @param request body RegisterReq true "request body"
// @produce json
// @success 200 {object} pkg.Response
// @router /api/user/register [post]
func register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	user := model.NewUser(req.Username, req.Password)
	err := user.Create()
	if err != nil {
		logger.Error(err.Error())
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "register fail", err.Error()))
		return
	}

	c.JSON(pkg.CreateSuccessResponse())
}
