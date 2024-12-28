package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *RedisClient
	NoDataFound = errors.New("No Data Found")
)

// RedisClient struct
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient get the redis client
func NewRedisClient(ctx context.Context, hostname string, password string) (CacheBase, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     hostname,
		Password: password,
		DB:       0, //Use default DB
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	redisClient = &RedisClient{Client: rdb}
	return redisClient, nil
}

// SetKey set the key
func (c *RedisClient) SetKey(ctx context.Context, key, value string) (err error) {
	err = c.Client.Set(ctx, key, value, 0).Err()
	if err != nil {
		err = fmt.Errorf("Redis Set Key Error: " + err.Error())
	}
	return
}

// Ping checks the status of the Redis server
func (c *RedisClient) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

// SetExpiringKey set the key with an expiration time
func (c *RedisClient) SetExpiringKey(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	err = c.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		err = fmt.Errorf("Redis Set Key Error: " + err.Error())
	}
	return
}

// GetKey get the key
func (c *RedisClient) GetKey(ctx context.Context, key string) (val string, err error) {
	val, err = c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", NoDataFound
	}
	if err != nil {
		return "", err
	}
	return
}

func (c *RedisClient) GetKeyAndFloatValue(ctx context.Context, key string) (val float64, err error) {
	val, err = c.Client.Get(ctx, key).Float64()
	if err == redis.Nil {
		return 0.0, NoDataFound
	}
	if err != nil {
		return 0.0, err
	}
	return
}

// IncrementKey inrement the given key
func (c *RedisClient) IncrementKey(ctx context.Context, key string) (err error) {
	err = c.Client.Incr(ctx, key).Err()
	if err != nil {
		err = fmt.Errorf("Redis Inrement Key Error: " + err.Error())
	}
	return
}

func (c *RedisClient) HMGet(ctx context.Context, key string, fields ...string) (fieldKeys []interface{}, err error) {
	fieldKeys, err = c.Client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		err = fmt.Errorf("Redis HMGet Error: " + err.Error())
	}
	return fieldKeys, err
}

// IncrementKeyValueBy increment the given key data by value
func (c *RedisClient) IncrementKeyValueBy(ctx context.Context, key string, value int64) (err error) {
	err = c.Client.IncrBy(ctx, key, value).Err()
	if err != nil {
		err = fmt.Errorf("Redis Inrement Key by value Error: " + err.Error())
	}
	return
}

// GetKeyUnMarshal get key and  unmarshal
func (c *RedisClient) GetKeyUnMarshal(ctx context.Context, key string, src interface{}) error {
	stringifiedData, err := c.GetKey(ctx, key)
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
func (c *RedisClient) SetExpiringKeyMarshal(ctx context.Context,
	key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.SetExpiringKey(ctx, key, string(cacheEntry), expiration)
	if err != nil {
		return err
	}
	return nil
}

// DeleteKey delete a key
func (c *RedisClient) DeleteKey(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisClient) DeleteMultipleKeys(ctx context.Context, keys []string) error {
	return c.Client.Del(ctx, keys...).Err()
}

// FetchAddresses fetch the addresses
func (c *RedisClient) FetchAddresses(ctx context.Context) map[string]string {
	addrMap := make(map[string]string)
	return addrMap
}

func (c *RedisClient) GetHashKey(ctx context.Context, hash string, key string) (val string, err error) {
	val, err = c.Client.HGet(ctx, hash, key).Result()
	if err == redis.Nil {
		return "", NoDataFound
	}
	if err != nil {
		return "", err
	}
	return
}

func (c *RedisClient) GetHashAll(ctx context.Context, key string) (val map[string]string, err error) {
	val, err = c.Client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, NoDataFound
	}
	if err != nil {
		return nil, err
	}
	return
}

func (c *RedisClient) DelHashKey(ctx context.Context, key string, field string) error {
	return c.Client.HDel(ctx, key, field).Err()
}

func (c *RedisClient) SetHashKey(ctx context.Context,
	hash string, key string, value interface{}, expiration time.Duration) (val string, err error) {
	err = c.Client.HSet(ctx, hash, key, value).Err()
	if err != nil {
		err = fmt.Errorf("Redis HSet Key Error: " + err.Error())
	}
	if expiration != 0 {
		err = c.Client.Expire(ctx, hash, expiration).Err()
		if err != nil {
			err = fmt.Errorf("Redis HSet Key Error: " + err.Error())
		}
	}

	return
}

// IncrementHKey increments the given key
func (c *RedisClient) IncrementHKey(ctx context.Context, hash string, key string, incrByValue ...int64) (val int64, err error) {
	// Set default value of incrBy to 1 if not provided
	incrBy := int64(1)
	if len(incrByValue) > 0 {
		incrBy = incrByValue[0]
	}

	// Increment the hash key using HIncrBy
	val, err = c.Client.HIncrBy(ctx, hash, key, incrBy).Result()
	if err != nil {
		err = fmt.Errorf("Redis Hash Increment Key Error: " + err.Error())
		return 0, err
	}

	return val, nil
}

// ExpireKey expires the given key after the given `expiration`
func (c *RedisClient) ExpireKey(ctx context.Context, key string, expiration time.Duration) (err error) {
	err = c.Client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("Redis expire Key Error: " + err.Error())
	}
	return nil
}

func (c *RedisClient) ExpireAt(ctx context.Context, key string, time time.Time) (err error) {
	err = c.Client.ExpireAt(ctx, key, time).Err()
	if err != nil {
		return fmt.Errorf("redis error setting expiry, err: %+v", err)
	}
	return nil
}

func (c *RedisClient) Subscribe(ctx context.Context, channels ...string) (pubsub *redis.PubSub, err error) {
	sub := c.Client.Subscribe(ctx, channels...)
	_, err = sub.Receive(ctx)
	if err != nil {
		err = fmt.Errorf("subscribe Key Error: " + err.Error())
		return &redis.PubSub{}, err
	}
	return sub, nil
}

func (c *RedisClient) LPush(ctx context.Context, key string, value interface{}) (val int64, err error) {
	val, err = c.Client.LPush(ctx, key, value).Result()
	if err != nil {
		err = fmt.Errorf("Redis LPush Error: " + err.Error())
	}
	return
}

func (c *RedisClient) RPop(ctx context.Context, key string) (val string, err error) {
	val, err = c.Client.RPop(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("Redis RPop Error: " + err.Error())
	}
	return
}

func (c *RedisClient) SAdd(ctx context.Context, key string, members ...string) (val int64, err error) {
	val, err = c.Client.SAdd(ctx, key, members).Result()
	if err != nil {
		err = fmt.Errorf("Redis sadd error: " + err.Error())
	}
	return
}

func (c *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) (val int64, err error) {
	val, err = c.Client.SRem(ctx, key, members).Result()
	if err != nil {
		err = fmt.Errorf("Redis remove error: " + err.Error())
	}
	return
}

func (c *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (val bool, err error) {
	val, err = c.Client.SIsMember(ctx, key, member).Result()
	if err != nil {
		err = fmt.Errorf("Redis SIsmember error: " + err.Error())
	}
	return
}

func (c *RedisClient) SScan(ctx context.Context, key string, crsr uint64, match string, count int64) (
	keys []string, cursor uint64, err error) {

	keys, cursor, err = c.Client.SScan(ctx, key, crsr, match, count).Result()
	if err != nil {
		err = fmt.Errorf("Redis SScan error: " + err.Error())
	}
	return
}

func (c *RedisClient) SPop(ctx context.Context, key string) (val string, err error) {
	val, err = c.Client.SPop(ctx, key).Result()

	if err != nil {
		err = fmt.Errorf("Redis SPop error: " + err.Error())
	}
	return
}

func (c *RedisClient) HKeys(ctx context.Context, key string) (keys []string, err error) {
	keys, err = c.Client.HKeys(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("Redis HKEYS error: " + err.Error())
	}
	return
}

func (c *RedisClient) ZAdd(ctx context.Context, key string, keyValue redis.Z) (err error) {
	result := c.Client.ZAdd(ctx, key, keyValue)
	if result.Err() != nil {
		err = fmt.Errorf("Redis ZAdd error: " + result.Err().Error())
	}
	return
}

func (c *RedisClient) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	result := c.Client.ZRevRange(ctx, key, start, stop)
	if result.Err() != nil {
		err := fmt.Errorf("Redis ZRevRange error: " + result.Err().Error())
		return nil, err
	}

	return result.Val(), nil
}

func (c *RedisClient) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	result := c.Client.ZRevRangeByScore(ctx, key, opt)
	if result.Err() != nil {
		err := fmt.Errorf("Redis ZRevRangeByScore error: " + result.Err().Error())
		return nil, err
	}

	return result.Val(), nil
}

func (c *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	result := c.Client.ZRange(ctx, key, start, stop)
	if result.Err() != nil {
		err := fmt.Errorf("Redis ZRange error: " + result.Err().Error())
		return nil, err
	}

	return result.Val(), nil
}

func (c *RedisClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	result := c.Client.ZRangeByScore(ctx, key, opt)
	if result.Err() != nil {
		err := fmt.Errorf("Redis ZRangeByScore error: " + result.Err().Error())
		return nil, err
	}

	return result.Val(), nil
}

func (c *RedisClient) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	result := c.Client.ZRemRangeByScore(ctx, key, min, max)
	if result.Err() != nil {
		err := fmt.Errorf("Redis ZRemRangeByScore error: " + result.Err().Error())
		return 0, err
	}

	return result.Val(), nil
}

// Keys will return matching pattern ex: [key:one key:three key:two] , try to avoid using this request may be expensive
func (c *RedisClient) Keys(ctx context.Context, keyPattern string) (keys []string, err error) {
	pattern := strings.TrimSpace(keyPattern)
	first := string(pattern[0])
	chungs := strings.Split(pattern, ":")
	if pattern == "" || first == "*" || len(chungs[0]) < 5 {
		err = fmt.Errorf("request is to expensive , pattern should not contain '*' OR less then five char in first string pattern of key")
	}
	keys, err = c.Client.Keys(ctx, keyPattern).Result()
	if err != nil {
		err = fmt.Errorf("Redis keys error: " + err.Error())
	}
	return
}

// SetNXExpiringKey set the key with an expiration time
func (c *RedisClient) SetNXExpiringKey(ctx context.Context, key, value string, expiration time.Duration) (err error) {
	err = c.Client.SetNX(ctx, key, value, expiration).Err()
	if err != nil {
		err = fmt.Errorf("Redis SetNX Key Error: " + err.Error())
	}
	return
}

func (c *RedisClient) FunctionLoad(ctx context.Context, script string) (*redis.StringCmd, error) {
	loadCmd := c.Client.FunctionLoad(ctx, script)
	if loadCmd.Err() != nil {
		return nil, loadCmd.Err()
	}
	return loadCmd, nil
}

func (c *RedisClient) FunctionCall(ctx context.Context, functionName string, keys []string, args ...interface{}) (*redis.Cmd, error) {
	loadCmd := c.Client.FCall(ctx, functionName, keys, args)
	if loadCmd.Err() != nil {
		return nil, loadCmd.Err()
	}
	return loadCmd, nil
}

func (c *RedisClient) FunctionReload(ctx context.Context, libName string, script string) (*redis.StringCmd, error) {
	loadCmd1 := c.Client.FunctionDelete(ctx, libName)
	if loadCmd1.Err() != nil {
		return nil, loadCmd1.Err()
	}
	loadCmd2 := c.Client.FunctionLoad(ctx, script)
	if loadCmd2.Err() != nil {
		return nil, loadCmd2.Err()
	}
	return loadCmd2, nil
}

func (c *RedisClient) KeyExists(ctx context.Context, key string) (exists bool, err error) {
	count, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("redis exists error: %v", err)
		return
	}
	exists = count > 0
	return
}
