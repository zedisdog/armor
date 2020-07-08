package helper

import jsoniter "github.com/json-iterator/go"

// Only pick items in data where key in keys
func Only(data map[string]interface{}, keys []string) (newMap map[string]interface{}) {
	newMap = make(map[string]interface{})
	for _, key := range keys {
		newMap[key] = data[key]
	}
	return
}

// Struct2Map covert struct to map
func Struct2Map(s interface{}) (m map[string]interface{}) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	dataJSON, _ := json.Marshal(s)
	json.Unmarshal(dataJSON, &m)
	return
}
