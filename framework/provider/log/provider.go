package log

import (
	"io"

	"github.com/skeleton1231/skeleton/framework"
	"github.com/skeleton1231/skeleton/framework/contract"
)

type HadeLogServiceProvider struct {
	framework.ServiceProvider

	Driver string

	// log Level
	Level contract.LogLevel
	// Log format
	Formatter contract.Formatter
	// Log of context
	CtxFielder contract.CtxFielder
	// Log Output
	Output io.Writer
}

// Register
func (l *HadeLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if l.Driver == "" {
		//tcs, err := c.Make(contract.AppKey)

	}
	return nil
}

// Boot 启动的时候注入
func (l *HadeLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer 是否延迟加载
func (l *HadeLogServiceProvider) IsDefer() bool {
	return false
}

func (l *HadeLogServiceProvider) Params(c framework.Container) []interface{} {
	return nil
}

// Name 定义对应的服务字符串凭证
func (l *HadeLogServiceProvider) Name() string {
	return contract.LogKey
}
