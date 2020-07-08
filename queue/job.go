package queue

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/zedisdog/armor/helper"
	"github.com/zedisdog/armor/log"
	"reflect"
)

var types = make(map[string]reflect.Type)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Job 队列任务接口
type Job interface {
	// Handle 执行任务时调用该方法
	Handle(c *Context)
}

// Register 注册job
func Register(job Job) {
	t := reflect.TypeOf(reflect.ValueOf(job).Elem().Interface())
	types[t.Name()] = t
}

// covert struct to json
func jobToJSON(job Job) (result []byte) {
	m := helper.Struct2Map(job)
	m["type_name"] = reflect.TypeOf(job).Elem().Name()
	result, _ = json.Marshal(m)
	return
}

func jsonToJob(jobJson []byte) Job {
	var m map[string]interface{}
	err := json.Unmarshal(jobJson, &m)
	if err != nil {
		log.Log.WithError(err).Warn("job parse failed")
	}

	t := types[m["type_name"].(string)]
	job := reflect.New(t).Interface()
	_ = json.Unmarshal(jobJson, job)
	return job.(Job)
}
