package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const redisRequestTimeout = time.Second * 5

func getRedisWaiter() waiter {
	host := envVar("REDIS_HOST", "localhost")
	port := envVar("REDIS_PORT", "6379")

	return redisWaiter{
		connectionString: fmt.Sprintf("%s:%s", host, port),
	}
}

type redisWaiter struct {
	connectionString string
}

func (w redisWaiter) waitFor() (bool, error) {
	client := redis.NewClient(&redis.Options{
		Addr: w.connectionString,
	})
	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), redisRequestTimeout)
	defer cancelFunc()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return false, err
	}

	s := strings.TrimSpace(pong)
	s = strings.ToLower(s)

	if s != "pong" {
		return false, errors.Errorf("%s ping-pong reponse '%s'", w.name(), s)
	}

	return true, nil
}

func (w redisWaiter) name() string {
	return "Redis"
}
