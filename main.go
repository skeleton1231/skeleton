package main

import (
	"github.com/skeleton1231/skeleton/app/console"
	"github.com/skeleton1231/skeleton/app/http"
	"github.com/skeleton1231/skeleton/framework"
	"github.com/skeleton1231/skeleton/framework/provider/app"
	"github.com/skeleton1231/skeleton/framework/provider/distributed"
	"github.com/skeleton1231/skeleton/framework/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewHadeContainer()
	// 绑定App服务提供者
	container.Bind(&app.HadeAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&distributed.LocalDistributedProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}

/*func main() {
	//core := framework.NewCore()
	core := gin.New()

	//core.Bind(&app.HadeAppProvider{})
	container := framework.NewHadeContainer()

	// 绑定 App 服务提供者
	container.Bind(&app.HadeAppProvider{})

	// 后续初始化需要绑定的服务提供者...
	// 将 HTTP 引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
	}

	//bind
	core.Bind(&demo.DemoServiceProvider{})
	//middleware
	core.Use(gin.Recovery())
	core.Use(middleware.Cost())

	registerRouter(core)
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: core,
		// 请求监听地址
		Addr: ":8881",
	}
	server.ListenAndServe()

	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}*/
