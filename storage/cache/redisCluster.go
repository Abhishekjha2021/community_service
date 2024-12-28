package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClusterClient *RedisClusterClient
)

// RedisClusterClient struct
type RedisClusterClient struct {
	Client *redis.ClusterClient
}

// NewRedisClusterClient redis clusterclient
func NewRedisClusterClient(ctx context.Context, hostname string, password string) (*RedisClusterClient, error) {
	hostnames := strings.Split(hostname, ",")
	var addr []string
	for i := 0; i < len(hostnames); i++ {
		addr = append(addr, hostnames[i])
	}
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addr,
		Password: password,
	})

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	redisClusterClient = &RedisClusterClient{Client: c}
	return redisClusterClient, nil
}

// Ping checks the status of the Redis server
func (r *RedisClusterClient) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// SetKey set the key
func (r *RedisClusterClient) SetKey(ctx context.Context, key, value string) (err error) {
	err = r.Client.Set(ctx, key, value, 0).Err()
	if err != nil {
		err = fmt.Errorf("Redis Set Key Error: " + err.Error())
	}
	return
}

// SetExpiringKey set the key with an expiration time
func (r *RedisClusterClient) SetExpiringKey(ctx context.Context,
	key, value string, expiration time.Duration) (err error) {
	err = r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		err = fmt.Errorf("Redis Set Key Error: " + err.Error())
	}
	return
}

// GetKey get the key
func (r *RedisClusterClient) GetKey(ctx context.Context, key string) (val string, err error) {
	val, err = r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return
}

// IncrementKey inrement the given key
func (r *RedisClusterClient) IncrementKey(ctx context.Context, key string) (err error) {
	err = r.Client.Incr(ctx, key).Err()
	if err != nil {
		err = fmt.Errorf("Redis Increment Key Error: " + err.Error())
	}
	return
}

// IncrementKeyValueBy inrement the given key data by value
func (c *RedisClusterClient) IncrementKeyValueBy(ctx context.Context, key string, value int64) (err error) {
	err = c.Client.IncrBy(ctx, key, value).Err()
	if err != nil {
		err = fmt.Errorf("Redis Inrement Key by value Error: " + err.Error())
	}
	return
}

// GetKeyUnMarshal get key and  unmarshal
func (r *RedisClusterClient) GetKeyUnMarshal(ctx context.Context, key string, src interface{}) error {
	stringifiedData, err := r.GetKey(ctx, key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(stringifiedData), &src)
	if err != nil {
		return err
	}
	return nil
}

// SetExpiringKeyMarshal get key and  unmarshal
func (r *RedisClusterClient) SetExpiringKeyMarshal(ctx context.Context,
	key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.SetExpiringKey(ctx, key, string(cacheEntry), expiration)
	if err != nil {
		return err
	}
	return nil
}

// DeleteKey delete a key
func (r *RedisClusterClient) DeleteKey(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// FetchAddresses fetch the addresses
func (r *RedisClusterClient) FetchAddresses() map[string]string {
	addrMap := make(map[string]string)
	return addrMap
}

func (c *RedisClusterClient) GetHashKey(ctx context.Context, hash string, key string) (val string, err error) {
	val, err = c.Client.HGet(ctx, hash, key).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return
}

func (c *RedisClusterClient) SetHashKey(ctx context.Context, hash string, key string, value interface{}, expiration time.Duration) (val string, err error) {
	err = c.Client.HSet(ctx, hash, key, value).Err()
	if err != nil {
		err = fmt.Errorf("Redis HSet Key Error: " + err.Error())
	}

	err = c.Client.Expire(ctx, hash, expiration).Err()
	if err != nil {
		err = fmt.Errorf("Redis HSet Key Error: " + err.Error())
	}

	return
}

// IncrementKey inrement the given key
func (c *RedisClusterClient) IncrementHKey(ctx context.Context, hash string, key string) (val int64, err error) {
	val, err = c.Client.HIncrBy(ctx, hash, key, 1).Result()
	if err != nil {
		err = fmt.Errorf("Redis Hash Inrement Key Error: " + err.Error())
		return 0, err
	}
	return val, nil
}

func (c *RedisClusterClient) ExpireKey(ctx context.Context, key string, expiration time.Duration) (err error) {
	err = c.Client.Expire(ctx, key, expiration).Err()
	if err != nil {
		err = fmt.Errorf("Redis expire Key Error: " + err.Error())
	}
	return
}
func (c *RedisClusterClient) Subscribe(ctx context.Context, channels ...string) (pubsub *redis.PubSub, err error) {
	sub := c.Client.Subscribe(ctx, channels...)
	_, err = sub.Receive(ctx)
	if err != nil {
		err = fmt.Errorf("subscribe Key Error: " + err.Error())
		return &redis.PubSub{}, err
	}
	return sub, nil
}

func (c *RedisClusterClient) LPush(ctx context.Context, key string, value interface{}) (val int64, err error) {
	val, err = c.Client.LPush(ctx, key, value).Result()
	if err != nil {
		err = fmt.Errorf("Redis LPush Error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) RPop(ctx context.Context, key string) (val string, err error) {
	val, err = c.Client.RPop(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("Redis RPop Error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) SAdd(ctx context.Context, key string, members ...string) (val int64, err error) {
	val, err = c.Client.SAdd(ctx, key, members).Result()
	if err != nil {
		err = fmt.Errorf("Redis sadd error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) SRem(ctx context.Context, key string, members ...interface{}) (val int64, err error) {
	val, err = c.Client.SRem(ctx, key, members).Result()
	if err != nil {
		err = fmt.Errorf("Redis remove error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) SIsMember(ctx context.Context, key string, member interface{}) (val bool, err error) {
	val, err = c.Client.SIsMember(ctx, key, member).Result()
	if err != nil {
		err = fmt.Errorf("Redis SIsmember error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) SPop(ctx context.Context, key string) (val string, err error) {
	val, err = c.Client.SPop(ctx, key).Result()

	if err != nil {
		err = fmt.Errorf("Redis SPop error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) SScan(ctx context.Context, key string, crsr uint64, match string, count int64) (
	keys []string, cursor uint64, err error) {

	keys, cursor, err = c.Client.SScan(ctx, key, crsr, match, count).Result()
	if err != nil {
		err = fmt.Errorf("Redis SScan error: " + err.Error())
	}
	return
}

func (c *RedisClusterClient) HKeys(ctx context.Context, key string) (keys []string, err error) {
	keys, err = c.Client.HKeys(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("Redis HKEYS error: " + err.Error())
	}
	return
}

// Keys will return matching pattern ex: [key:one key:three key:two] , try to avoid using this request may be expensive
func (c *RedisClusterClient) Keys(ctx context.Context, keyPattern string) (keys []string, err error) {
	pattern := strings.TrimSpace(keyPattern)
	first := string(pattern[0])
	chungs := strings.Split(pattern, ":")
	if pattern == "" || first == "*" || len(chungs[0]) < 5 {
		err = fmt.Errorf("request is to expensive , pattern should not contain '*' OR less then five char in first string pattern of key")
		return nil,err
	}
	keys, err = c.Client.Keys(ctx, keyPattern).Result()
	if err != nil {
		err = fmt.Errorf("Redis keys error: " + err.Error())
	}
	return
}
