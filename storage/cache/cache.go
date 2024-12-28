package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheBase tokengen interface
type CacheBase interface {
	SetKey(context.Context, string, string) error
	SetExpiringKey(context.Context, string, interface{}, time.Duration) error
	GetKey(context.Context, string) (string, error)
	GetKeyAndFloatValue(ctx context.Context, key string) (float64, error)
	IncrementKey(context.Context, string) error
	IncrementKeyValueBy(ctx context.Context, key string, value int64) (err error)
	Ping(context.Context) error
	GetKeyUnMarshal(context.Context, string, interface{}) error
	SetExpiringKeyMarshal(context.Context, string, interface{}, time.Duration) error
	DeleteKey(context.Context, string) error
	DeleteMultipleKeys(ctx context.Context, keys []string) error
	FetchAddresses(context.Context) map[string]string
	GetHashKey(ctx context.Context, hash string, key string) (val string, err error)
	GetHashAll(ctx context.Context, key string) (val map[string]string, err error)
	SetHashKey(ctc context.Context, hash string, key string, value interface{}, expiration time.Duration) (val string, err error)
	IncrementHKey(ctx context.Context, hash string, key string, incrByValue ...int64) (val int64, err error)
	ExpireKey(ctx context.Context, key string, expiration time.Duration) (err error)
	Subscribe(ctx context.Context, channels ...string) (pubsub *redis.PubSub, err error)
	LPush(ctx context.Context, key string, value interface{}) (val int64, err error)
	RPop(ctx context.Context, key string) (val string, err error)
	SAdd(ctx context.Context, key string, members ...string) (val int64, err error)
	SRem(ctx context.Context, key string, members ...interface{}) (val int64, err error)
	SIsMember(ctx context.Context, key string, member interface{}) (val bool, err error)
	SScan(ctx context.Context, key string, crsr uint64, match string, count int64) (keys []string, cursor uint64, err error)
	HKeys(ctx context.Context, key string) (keys []string, err error)
	HMGet(ctx context.Context, key string, fields ...string) (fieldKeys []interface{}, err error)
	SPop(ctx context.Context, key string) (val string, err error)
	DelHashKey(ctx context.Context, key string, field string) error
	ZAdd(ctx context.Context, key string, keyValue redis.Z) (err error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error)
	//Keys will return matching pattern ex: [key:one key:three key:two] , try to avoid using this request may be expensive
	Keys(ctx context.Context, keyPattern string) ([]string, error)
	SetNXExpiringKey(ctx context.Context, key, value string, expiration time.Duration) error
	ExpireAt(ctx context.Context, key string, time time.Time) (err error)
	FunctionLoad(ctx context.Context, script string) (*redis.StringCmd, error)
	FunctionCall(ctx context.Context, functionName string, keys []string, args ...interface{}) (*redis.Cmd, error)
	FunctionReload(ctx context.Context, libName string, script string) (*redis.StringCmd, error)
	KeyExists(ctx context.Context, key string) (exists bool, err error) 
}
