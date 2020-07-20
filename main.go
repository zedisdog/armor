package armor

import (
	"context"
	redis "github.com/zedisdog/armor/cache"
	"os"
	"os/signal"
	"sync"
	"syscall"
	casbin2 "github.com/zedisdog/armor/casbin"
	"github.com/zedisdog/armor/config"
	"github.com/zedisdog/armor/model"
	"github.com/zedisdog/armor/queue"
	"github.com/zedisdog/armor/web"
)

type armor struct {
	enableQueue   bool                 // 是否启用队列
	enableCache   bool                 // 是否使用缓存
	autoMigrate   model.AutoMigrate    // 自动迁移方法
	configPath    string               // 配置文件路劲
	routesMaker   *web.RoutesMaker     // 路由方法
	casbinOptions []casbin2.ConfigFunc // 访问控制配置
}

type ConfigFunc func(*armor)

func WithQueue(enabled bool) ConfigFunc {
	return func(a *armor) {
		a.enableQueue = enabled
	}
}

func WithCache(enabled bool) ConfigFunc {
	return func(a *armor) {
		a.enableCache = enabled
	}
}

func WithAutoMigrate(m model.AutoMigrate) ConfigFunc {
	return func(a *armor) {
		a.autoMigrate = m
	}
}

func WithCasbin(configs ...casbin2.ConfigFunc) ConfigFunc {
	return func(a *armor) {
		a.casbinOptions = configs
	}
}

func WithConfigPath(s string) ConfigFunc {
	return func(a *armor) {
		a.configPath = s
	}
}

func WithRoutesMaker(maker *web.RoutesMaker) ConfigFunc {
	return func(a *armor) {
		a.routesMaker = maker
	}
}

func NewArmor(configs ...ConfigFunc) *armor {
	a := &armor{}
	for _, conf := range configs {
		conf(a)
	}
	if a.configPath == "" {
		a.configPath = "config.yml" // 默认值
	}
	return a
}

func (a *armor) Start() error {
	config.Init(a.configPath)
	model.Init()
	if a.autoMigrate != nil {
		a.autoMigrate(model.DB)
	}
	if a.casbinOptions != nil {
		casbin2.Init(casbin2.NewOptions(a.casbinOptions...))
	}

	cxt, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if a.enableCache {
		cacheInstance, err := redis.InitCache()
		defer cacheInstance.Close()
		if err != nil {
			return err
		}
	}
	if a.enableQueue {
		queue.Init()
		err := queue.Instance().Start(cxt, &wg)
		defer queue.Instance().Close()
		if err != nil {
			return err
		}
	}
	web.Start(cxt, &wg, a.routesMaker)

	<-sigs
	cancel()
	wg.Wait()

	return nil
}
