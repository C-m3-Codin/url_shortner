package utils

import (
	"context"
	"encoding/json"

	"github.com/c-m3-codin/url_shortner/services"
)

func GetRedisValue(key string) (val string, err error) {
	//fmt.Println"Key for redis is \n\n", key)
	val, err = services.RedisClient.Get(context.Background(), key).Result()
	//fmt.Println"value for redis is \n\n", val)
	return
}

func SetRedisValue(key string, value interface{}) (result string, err error) {
	//fmt.Printlnvalue)
	jsonData, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	//fmt.PrintlnjsonData)

	// Replace "yourKey" with the key where you want to store the JSON data
	//fmt.Println"key is ", key)
	err = services.RedisClient.Set(context.Background(), key, jsonData, 0).Err()
	if err != nil {
		panic(err)
	}

	return

}
