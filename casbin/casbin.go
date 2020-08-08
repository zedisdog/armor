package casbin

import (
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"regexp"
)

type Options struct {
	PolicyFilePath string
	Adapter        persist.Adapter
	Model          casbinModel.Model
}

func NewOptions(v *viper.Viper, db *gorm.DB) *Options {
	o := new(Options)
	if v.IsSet("casbin.PolicyFilePath") && !v.IsSet("casbin.Adapter") {
		a, _ := gormadapter.NewAdapterByDB(db)
		o.Adapter = a
	}
	if !v.IsSet("casbin.Model") {
		m, _ := casbinModel.NewModelFromString(DEFAULT_RBAC_CONFIG)
		o.Model = m
	} else {
		o.Model = genModel(v.Get("casbin.Model"))
	}

	return o
}

func New(v *viper.Viper, options *Options) (Enforcer *casbin.Enforcer, err error) {
	if !v.IsSet("casbin.enable") || !v.GetBool("casbin.enable") {
		return nil, nil
	}
	if options.Adapter != nil { //通过adapter或者策略文件路径实例化
		Enforcer, err = casbin.NewEnforcer(options.Model, options.Adapter)
	} else {
		Enforcer, err = casbin.NewEnforcer(options.Model, options.PolicyFilePath)
	}
	return
}

func genModel(m interface{}) (cm casbinModel.Model) {
	switch m.(type) {
	case casbinModel.Model: // 判断是model对象就就直接赋值
		cm = m.(casbinModel.Model)
	case string: // 判断是字符串 就要判断是配置内容还是配置文件路径 再根据类型选择相应的实例化方法
		// 多行就表示是配置内容 单行就表示是配置文件路劲
		if isMultiLine(m.(string)) {
			cm, _ = casbinModel.NewModelFromString(m.(string))
		} else {
			cm, _ = casbinModel.NewModelFromFile(m.(string))
		}
	}
	return
}

func isMultiLine(s string) bool {
	// 有换行说明是多行
	matched, _ := regexp.MatchString("/[\r\n]+/", s)
	return matched
}

var ProviderSet = wire.NewSet(NewOptions, New)
