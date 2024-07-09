package user

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.POST("/login", login)
	r.POST("/register", register)
}
