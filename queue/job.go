package queue

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/zedisdog/armor/log"
	"reflect"
)

var types map[string]reflect.Type = make(map[string]reflect.Type)

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
func jobToJSON(data Job) []byte {
	var m map[string]interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	dataJSON, _ := json.Marshal(data)
	json.Unmarshal(dataJSON, &m)
	m["type_name"] = reflect.TypeOf(data).Elem().Name()
	newJSON, _ := json.Marshal(m)
	return newJSON
}

func jsonToJob(data []byte) Job {
	var m map[string]interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Log.WithError(err).Warn("job parse failed")
	}

	t := types[m["type_name"].(string)]
	job := reflect.New(t).Interface()
	json.Unmarshal(data, job)
	return job.(Job)
}
