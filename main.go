package main

import (
	"github.com/skeleton1231/skeleton/app/console"
	"github.com/skeleton1231/skeleton/app/http"
	"github.com/skeleton1231/skeleton/app/provider/demo"
	"github.com/skeleton1231/skeleton/framework"
	"github.com/skeleton1231/skeleton/framework/provider/app"
	"github.com/skeleton1231/skeleton/framework/provider/config"
	"github.com/skeleton1231/skeleton/framework/provider/distributed"
	"github.com/skeleton1231/skeleton/framework/provider/env"
	"github.com/skeleton1231/skeleton/framework/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewHadeContainer()
	// 绑定App服务提供者
	container.Bind(&app.HadeAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.HadeEnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&config.HadeConfigProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}
	container.Bind(&demo.DemoProvider{})

	// 运行root命令
	console.RunCommand(container)
}
