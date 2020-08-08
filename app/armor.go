package app

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/zedisdog/armor/model"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var app *Armor

type Armor struct {
	autoMigrate model.AutoMigrate  // 自动迁移方法
	Config      *viper.Viper       //配置文件
	Queue       Queue              // 队列
	Enforcer    *casbin.Enforcer   // 访问控制
	Cache       *redis.Client      // redis缓存
	DB          *gorm.DB           //数据库
	Routes      Routes             // 路由
	seeders     []func(*gorm.DB)   // 数据库填充方法组
	CancelCxt   context.Context    // 关闭上下文
	Cancel      context.CancelFunc // 关闭函数
	Wg          *sync.WaitGroup    // 等待
	HttpServer  HttpServer         // http服务器
}

func New(config *viper.Viper, queue Queue, enforcer *casbin.Enforcer, cache *redis.Client, db *gorm.DB, httpserver HttpServer) *Armor {
	if app == nil {
		cxt, cancel := context.WithCancel(context.Background())
		app = &Armor{
			Config:     config,
			Queue:      queue,
			Enforcer:   enforcer,
			Cache:      cache,
			DB:         db,
			CancelCxt:  cxt,
			Cancel:     cancel,
			Wg:         &sync.WaitGroup{},
			HttpServer: httpserver,
		}
	}
	return app
}

func (a *Armor) SetMigrate(m model.AutoMigrate) {
	a.autoMigrate = m
}

func (a *Armor) SetSeeders(seeders ...func(*gorm.DB)) {
	a.seeders = seeders
}

func (a *Armor) SetRoutes(routes Routes) {
	a.Routes = routes
}

func (a *Armor) Start() error {
	// 数据库迁移
	if a.autoMigrate != nil {
		a.autoMigrate(a.DB)
	}
	// 数据填充
	if len(a.seeders) > 0 {
		for _, seeder := range a.seeders {
			seeder(a.DB)
		}
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if a.Queue != nil {
		if err := a.Queue.Start(a); err != nil {
			panic(err)
		}
	}
	if a.HttpServer != nil {
		a.HttpServer.Start(a)
	}

	<-sigs
	a.Cancel()
	a.Wg.Wait()

	return nil
}

var ProviderSet = wire.NewSet(New)
