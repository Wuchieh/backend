package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
	_ "string_backend_0001/docs"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/pkg"
	"string_backend_0001/internal/web/google"
	"string_backend_0001/internal/web/line"
	"string_backend_0001/internal/web/user"
)

func Run() error {
	r := gin.Default()
	router(r)
	webC := conf.Conf.Web
	addr := fmt.Sprintf("%s:%d", webC.Host, webC.Port)
	return r.Run(addr)
}

func router(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	api.GET("/hello", helloWorld)

	// 若不想使用 google oauth2 請註釋掉下列兩行
	googleApi := api.Group("/google")
	google.Router(googleApi)

	lineApi := api.Group("/line")
	line.Router(lineApi)

	userApi := api.Group("/user")
	user.Router(userApi)
}

// @summary hello world
// @description response hello world
// @tags Api
// @id helloWorld
// @produce json
// @success 200 {object} pkg.Response
// @router /api/hello [get]
func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, pkg.CreateSuccessResponseObj("Hello World!"))
}
