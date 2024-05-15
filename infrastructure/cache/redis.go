package cache

/*
import (
	log "xxx-server/application/logger"
	"xxx-server/infrastructure/config"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var redisClient redis.UniversalClient

func InitRedis() error {
	log.Debug("setting up Redis client")

	opt, err := redis.ParseURL(config.C.Redis.Addr)
	if err != nil {
		return errors.Wrap(err, "parse redis url error")
	}
	opt.PoolSize = config.C.Redis.PoolSize
	if config.C.Redis.Cluster {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{opt.Addr},
			PoolSize: opt.PoolSize,
			Password: opt.Password,
		})
	} else if config.C.Redis.MasterName != "" {
		redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       config.C.Redis.MasterName,
			SentinelAddrs:    []string{opt.Addr},
			SentinelPassword: opt.Password,
			DB:               opt.DB,
			PoolSize:         opt.PoolSize,
		})
	} else {
		redisClient = redis.NewClient(opt)
	}
	return nil
}

// RedisClient returns the Redis client.
func RedisClient() redis.UniversalClient {
	return redisClient
}
*/
