//+build wireinject

package armor

import (
	"github.com/google/wire"
	"github.com/zedisdog/armor/app"
	"github.com/zedisdog/armor/cache"
	"github.com/zedisdog/armor/casbin"
	"github.com/zedisdog/armor/config"
	"github.com/zedisdog/armor/model"
	"github.com/zedisdog/armor/queue"
	"github.com/zedisdog/armor/web"
)

var ProviderSet = wire.NewSet(
	config.ProviderSet,
	model.ProviderSet,
	web.ProviderSet,
	queue.ProviderSet,
	casbin.ProviderSet,
	cache.ProviderSet,
	app.ProviderSet,
)

func InitApp(config string) (*app.Armor, error) {
	wire.Build(providerSet)
	return nil, nil
}
