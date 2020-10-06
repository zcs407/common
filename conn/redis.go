package conn

import (
	"time"

	"github.com/go-redis/redis/v7"
)

type RClient struct {
	Pool *redis.ClusterClient
}

// redis读写连接池
func NewRedisRW(dsn []string, user, psd string) (redisClient *RClient, err error) {
	redisClient = &RClient{}
	redisClient.Pool = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        dsn,
		Username:     user,
		Password:     psd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		MaxRetries:   2,
		IdleTimeout:  5 * time.Minute,
	})
	_, err = redisClient.Pool.Ping().Result()
	return
}

// redis只读连接池
func NewRedisRO(dsn []string, user, psd string) (redisClient *RClient, err error) {
	redisClient = &RClient{}
	redisClient.Pool = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       dsn,
		Username:    user,
		Password:    psd,
		ReadOnly:    true,
		DialTimeout: 10 * time.Second,
		ReadTimeout: 30 * time.Second,
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
		MaxRetries:  2,
		IdleTimeout: 5 * time.Minute,
	})

	_, err = redisClient.Pool.Ping().Result()
	return
}

func (r *RClient) Close() {
	_ = r.Pool.Close()
}

// 设置string类型
func (r *RClient) SetStr(k, v string, expiration int64) error {
	return r.Pool.Set(k, v, time.Duration(expiration)).Err()
}

// 获取string类型的值
func (r *RClient) GetStr(k string) (v string, err error) {
	return r.Pool.Get(k).Result()
}
