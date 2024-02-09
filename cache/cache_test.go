package cache

import (
	"context"
	"testing"
	"time"

	mock_log "github.com/nafisalfiani/ketson-go-lib/tests/mock/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func Test_redis_Get(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		want          interface{}
		wantErr       bool
	}{
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				key: "test1",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				key: "test1",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				db.Set(context.Background(), "test1", "test1", time.Hour)

				return db
			},
			want:    "test1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		})
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()
			c := cache{
				rdb: rdb,
			}
			got, err := c.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cache.Get() = %v, want %v", got, tt.want)
			}
			db.Del(context.Background(), "test1")
		})
	}
}

func Test_redis_SetEX(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTime := time.Hour * 24

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	type args struct {
		ctx     context.Context
		key     string
		val     string
		expTime time.Duration
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		wantErr       bool
	}{
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				key:     "testset",
				val:     "yes",
				expTime: 0,
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6378",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				key:     "whatever",
				val:     "yes",
				expTime: mockTime,
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: false,
		},
		{
			name: "success with default ttl",
			args: args{
				ctx: context.Background(),
				key: "whatever",
				val: "yes",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		})
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()
			c := cache{
				conf: Config{
					DefaultTTL: mockTime,
				},
				rdb: rdb,
			}
			if err := c.SetEX(tt.args.ctx, tt.args.key, tt.args.val, tt.args.expTime); (err != nil) != tt.wantErr {
				t.Errorf("cache.SetEX() error = %v, wantErr %v", err, tt.wantErr)
			}
			db.Del(context.Background(), "testset")
		})
	}
}

func Test_cache_Del(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		prepIterMock  func() string
		wantErr       bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				key: "skey1",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			prepIterMock: func() string {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				db.Set(context.Background(), "skey1", "key1val", time.Hour)
				x := db.Scan(context.Background(), 0, "skey1", 0).Iterator()
				x.Next(context.Background())
				y := x.Val()
				return y
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()
			c := cache{
				rdb: rdb,
				log: logger,
			}
			str := tt.prepIterMock()
			if err := c.Del(tt.args.ctx, str); (err != nil) != tt.wantErr {
				t.Errorf("cache.Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
