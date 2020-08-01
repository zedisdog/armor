package casbin

import (
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/zedisdog/armor/model"
	"regexp"
)

var Enforcer *casbin.Enforcer

type Options struct {
	PolicyFilePath string
	Adapter        persist.Adapter
	Model          casbinModel.Model
}

func Init(options *Options) {
	var err error
	if options.Adapter != nil { //通过adapter或者策略文件路径实例化
		Enforcer, err = casbin.NewEnforcer(options.Model, options.Adapter)
	} else {
		Enforcer, err = casbin.NewEnforcer(options.Model, options.PolicyFilePath)
	}
	if err != nil {
		panic(err)
	}
}

func GetEnforcer(options *Options) (c *casbin.Enforcer) {
	if Enforcer == nil {
		Init(options)
	}
	return Enforcer
}

func NewOptions(configFuncs ...ConfigFunc) *Options {
	options := &Options{}
	for _, config := range configFuncs {
		config(options)
	}

	if options.PolicyFilePath == "" && options.Adapter == nil {
		a, _ := gormadapter.NewAdapterByDB(model.DB)
		options.Adapter = a
	}

	if options.Model == nil {
		m, _ := casbinModel.NewModelFromString(DEFAULT_RBAC_CONFIG)
		options.Model = m
	}

	return options
}

type ConfigFunc func(options *Options)

func WithPolicyFilePath(path string) ConfigFunc {
	return func(options *Options) {
		options.PolicyFilePath = path
		options.Adapter = nil
	}
}

func WithAdapter(a persist.Adapter) ConfigFunc {
	return func(options *Options) {
		options.Adapter = a
	}
}

func WithModel(m interface{}) ConfigFunc {
	cm := genModel(m)
	return func(options *Options) {
		options.Model = cm
	}
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
