package http

import (
	"github.com/skeleton1231/skeleton/app/http/module/demo"
	"github.com/skeleton1231/skeleton/framework/gin"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
