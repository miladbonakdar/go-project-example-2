package logic

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"hotel-engine/core"
	"hotel-engine/infrastructure/config"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

type redisDistributedLock struct {
	locker *redislock.Client
}

func (s *redisDistributedLock) Lock(key string, duration time.Duration, toDo func()) error {
	lock, err := s.locker.Obtain(ctx, key, duration, nil)
	if err == redislock.ErrNotObtained {
		return nil
	} else if err != nil {
		return err
	}
	defer lock.Release(ctx)
	toDo()
	return nil
}

type fakeDevelopmentLocker struct {
}

func (s *fakeDevelopmentLocker) Lock(key string, duration time.Duration, toDo func()) error {
	toDo()
	return nil
}

func NewRedisLocker() core.DistributedLocker {
	con := config.Get()
	if con.IsDevelopment() {
		return &fakeDevelopmentLocker{}
	}
	connectionDetails := strings.Split(con.MemoryStorageConnection, ",")
	db, err := strconv.Atoi(connectionDetails[2])
	if err != nil {
		db = 0
	}
	locker := redislock.New(redis.NewClient(&redis.Options{
		Addr:     connectionDetails[0],
		Password: connectionDetails[1],
		DB:       db,
	}))
	return &redisDistributedLock{
		locker: locker,
	}
}
