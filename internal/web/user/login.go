package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"string_backend_0001/internal/logger"
	"string_backend_0001/internal/model"
	"string_backend_0001/internal/pkg"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @summary user login
// @description user login
// @tags User
// @id login
// @accept json
// @param request body LoginReq true "request body"
// @produce json
// @success 200 {object} pkg.Response
// @router /api/user/login [post]
func login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		return
	}

	var err error
	user := model.NewUser(req.Username, req.Password)
	user, err = model.GetUser(*user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "login fail", err.Error()))
		return
	}

	c.JSON(pkg.CreateSuccessResponse(user.ID))
}
