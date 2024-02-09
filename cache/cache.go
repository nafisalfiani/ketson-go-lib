package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/nafisalfiani/ketson-go-lib/codes"
	"github.com/nafisalfiani/ketson-go-lib/errors"
	"github.com/nafisalfiani/ketson-go-lib/log"
	"github.com/redis/go-redis/v9"
)

var ErrNotObtained = redislock.ErrNotObtained

const (
	Nil = redis.Nil
)

type Locker *redislock.Lock

type Interface interface {
	Get(ctx context.Context, key string) (string, error)
	SetEX(ctx context.Context, key string, val string, expTime time.Duration) error
	Del(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
	FlushAllAsync(ctx context.Context) error
	FlushDB(ctx context.Context) error
	FlushDBAsync(ctx context.Context) error
	GetDefaultTTL(ctx context.Context) time.Duration
}

type TLSConfig struct {
	Enabled            bool `env:"ENABLED"`
	InsecureSkipVerify bool `env:"INSECURE_SKIP_VERIFY"`
}

type Config struct {
	Protocol   string        `env:"PROTOCOL"`
	Host       string        `env:"HOST"`
	Port       string        `env:"PORT"`
	Username   string        `env:"USERNAME"`
	Password   string        `env:"PASSWORD"`
	DefaultTTL time.Duration `env:"DEFAULT_TTL"`
	TLS        TLSConfig     `env:"TLS"`
}

type cache struct {
	conf  Config
	rdb   *redis.Client
	log   log.Interface
	rlock *redislock.Client
}

func Init(cfg Config, log log.Interface) Interface {
	log.Info(context.Background(), "connecting to redis...")

	c := &cache{
		conf: cfg,
		log:  log,
	}
	c.connect(context.Background())
	return c
}

func (c *cache) connect(ctx context.Context) {
	redisOpts := redis.Options{
		Network:  c.conf.Protocol,
		Addr:     fmt.Sprintf("%s:%s", c.conf.Host, c.conf.Port),
		Username: c.conf.Username,
		Password: c.conf.Password,
	}

	if c.conf.TLS.Enabled {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.conf.TLS.InsecureSkipVerify,
		}
	}

	client := redis.NewClient(&redisOpts)

	err := client.Ping(ctx).Err()
	if err != nil {
		c.log.Fatal(ctx, fmt.Sprintf("[FATAL] cannot connect to redis on address @%s:%v, with error: %s", c.conf.Host, c.conf.Port, err))
	}
	c.rdb = client
	c.log.Info(ctx, fmt.Sprintf("REDIS: Address @%s:%v", c.conf.Host, c.conf.Port))

	c.rlock = redislock.New(client)
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	s, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (c *cache) SetEX(ctx context.Context, key string, val string, expTime time.Duration) error {

	if expTime <= 0 {
		expTime = c.conf.DefaultTTL
	}

	err := c.rdb.SetEx(ctx, key, val, expTime).Err()
	if err != nil {
		return errors.NewWithCode(codes.CodeRedisSetex, err.Error())
	}

	return nil
}

func (c *cache) Del(ctx context.Context, key string) error {
	var keysCount int64
	// Use SCAN with COUNT = 0 to advance the cursor
	iter := c.rdb.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		c.log.Info(ctx, fmt.Sprintf("deleted key: %s", iter.Val()))
		c.rdb.Del(ctx, iter.Val())
		keysCount++
	}
	if err := iter.Err(); err != nil {
		return err
	}
	c.log.Info(ctx, fmt.Sprintf("sucessfuly deleted %d numbers of key", keysCount))

	return nil
}

func (c *cache) FlushAll(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}

func (c *cache) FlushAllAsync(ctx context.Context) error {
	return c.rdb.FlushAllAsync(ctx).Err()
}

func (c *cache) FlushDB(ctx context.Context) error {
	return c.rdb.FlushDB(ctx).Err()
}

func (c *cache) FlushDBAsync(ctx context.Context) error {
	return c.rdb.FlushDBAsync(ctx).Err()
}

func (c *cache) GetDefaultTTL(ctx context.Context) time.Duration {
	return c.conf.DefaultTTL
}
