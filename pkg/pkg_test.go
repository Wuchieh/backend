package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var c *gin.Context

func ExampleGetGoogleUserinfo() {
	u, _ := GetGoogleUserinfo(c)
	fmt.Println(u)
}
