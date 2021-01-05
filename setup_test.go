package rediscript

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest/v3"
)

var (
	dockerRes *dockertest.Resource
	redisPool *redis.Pool
	TTL       = int64(2)
)

func TestMain(m *testing.M) {
	setup()

	defer dockerRes.Close()

	m.Run()
}

func setup() {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("could not connect to docker, " + err.Error())
	}
	dockerRes, err = dockerPool.Run("redis", "5.0", nil)
	if err != nil {
		log.Fatal("could not start resource, " + err.Error())
	}

	redisPool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(fmt.Sprintf("redis://localhost:%s", dockerRes.GetPort("6379/tcp")))
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}
