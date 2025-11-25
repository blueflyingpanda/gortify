package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	rdbNum, err := strconv.Atoi(Conf.RedisName)
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", Conf.RedisHost, Conf.RedisPort), Password: Conf.RedisPass, DB: rdbNum})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}
