package main

import (
	"time"

	"github.com/skeleton1231/skeleton/framework/gin"
)

func UserLoginController(c *gin.Context) {
	foo, _ := c.DefaultQueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second)
	// 输出结果
	c.ISetOkStatus().IJson("ok, UserLoginController: " + foo)

}
