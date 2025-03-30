package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spigcoder/LittleBook/webook/interanal/domain"
)

var KeyNotExist = redis.Nil

type UserCache struct {
	client     redis.Cmdable
	expireTime time.Duration
}

// 面向接口编程，依赖注入，这里的cmd可以是本地的redis，也可以是集群的redis
// 这里的cmd是一个接口，它的实现是redis.Client和redis.ClusterClient
func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		client:     cmd,
		expireTime: time.Minute * 15,
	}
}

func (cache *UserCache) GetById(ctx context.Context, id int64) (domain.User, error) {
	key := cache.getKey(id)
	val, err := cache.client.Get(ctx, key).Bytes()	
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	//redis设置数据不能直接传递结构体，要先将其序列化，然后再进行缓存存储
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.getKey(u.Id)
	return cache.client.Set(ctx, key, val, cache.expireTime).Err()
}

func (cache *UserCache) getKey(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}
