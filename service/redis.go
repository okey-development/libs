package service

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var (
	once       sync.Once
	instance   *redis.Client
	subscriber func(payload string)
	muRedis    = sync.Mutex{}
	chRedis    <-chan *redis.Message
)

func client() *redis.Client {
	once.Do(func() {
		instance = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Адрес и порт Redis сервера
			Password: "",               // Пароль, если требуется
			DB:       0,                // Индекс базы данных Redis
		})
	})

	return instance
}

func SendRedisMessage(recipient string, message interface{}) error {
	err := client().Publish(recipient, message).Err()
	if err != nil {
		return Errorf(err.Error())
	}

	return nil
}

func Subscribe(callBack func(payload string), amountRoutines int) error {
	if callBack == nil {
		return Errorf("callBack == nil")
	}

	muRedis.Lock()
	if subscriber != nil {
		return Errorf("subscription already initialized")
	}
	subscriber = callBack
	sub := client().Subscribe(appConfig.AppName)
	chRedis = sub.Channel()
	muRedis.Unlock()

	for i := 0; i < amountRoutines; i++ {
		go func() {
			for msg := range chRedis {
				if msg.Channel == appConfig.AppName {
					subscriber(msg.Payload)
				}
			}
		}()
	}
	return nil
}

func SetStringValue(key, value string, expiration time.Duration) error {
	return client().Set(key, value, expiration).Err()
}

func SetInt64Value(key string, value int64, expiration time.Duration) error {
	return client().Set(key, value, expiration).Err()
}

func GetStringValue(key string) (string, error) {
	return client().Get("key").Result()
}

func GetStringInt64(key string) (int64, error) {
	return client().Get(key).Int64()
}
