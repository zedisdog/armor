package armor

import (
	"context"
	redis "github.com/zedisdog/armor/cache"
	casbin2 "github.com/zedisdog/armor/casbin"
	"github.com/zedisdog/armor/config"
	"github.com/zedisdog/armor/model"
	"github.com/zedisdog/armor/queue"
	"github.com/zedisdog/armor/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type armor struct {
	enableQueue   bool                 // 是否启用队列
	enableCache   bool                 // 是否使用缓存
	autoMigrate   model.AutoMigrate    // 自动迁移方法
	configPath    string               // 配置文件路劲
	routes        web.Routes           // 路由
	casbinOptions []casbin2.ConfigFunc // 访问控制配置
	seeders       []func()             // 数据库填充方法组
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

func WithRoutes(routes web.Routes) ConfigFunc {
	return func(a *armor) {
		a.routes = routes
	}
}

func WithSeeder(funcs ...func()) ConfigFunc {
	return func(a *armor) {
		a.seeders = funcs
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
	// 初始化配置
	config.Init(a.configPath)
	// 初始化数据库
	model.Init()
	// 数据库迁移
	if a.autoMigrate != nil {
		a.autoMigrate(model.DB)
	}
	// 数据填充
	if len(a.seeders) > 0 {
		for _, seeder := range a.seeders {
			seeder()
		}
	}
	// 初始化访问控制组件
	if a.casbinOptions != nil {
		casbin2.Init(casbin2.NewOptions(a.casbinOptions...))
	}
	// 初始化redis库
	if a.enableCache {
		cacheInstance, err := redis.InitCache()
		defer cacheInstance.Close()
		if err != nil {
			return err
		}
	}

	cxt, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// 开启队列
	if a.enableQueue {
		queue.Init()
		err := queue.Instance().Start(cxt, &wg)
		defer queue.Instance().Close()
		if err != nil {
			return err
		}
	}
	// 开始web服务
	web.Start(cxt, &wg, a.routes)

	<-sigs
	cancel()
	wg.Wait()

	return nil
}
