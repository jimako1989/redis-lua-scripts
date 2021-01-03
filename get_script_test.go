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
		MaxIdle:     3,
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

func TestScriptTTLAT(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	TTL := int64(100)

	script_ttlat, err := GetScript("util/1_ttlat.lua")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	conn.Do("SETEX", "Alice", TTL, "Bob")

	now := time.Now().Unix()
	expireAt, err := redis.Int64(script_ttlat.Do(conn, "Alice"))
	if err != nil {
		t.Fatalf("error while doing script, %v", err.Error())
	}

	if expireAt != now+TTL {
		t.Fatalf("error, expireAt: %v, time+TTL: %v", expireAt, now+TTL)
	}
}

// func TestScriptHSETEX(t *testing.T) {
// 	conn := redisPool.Get()
// 	defer conn.Close()

// 	TTL := 100

// 	script_hsetex, err := GetScript("lua/hashes_xp/1_hsetxp.lua")
// 	if err != nil {
// 		t.Fatalf("error connection to script, %v", err)
// 	}

// 	_, err = script_hsetex.Do(conn, "key", "field", "value", TTL)
// 	if err != nil {
// 		t.Fatalf("error to hsetex, %v", err)
// 	}

// 	ttl, err := redis.Int(conn.Do("TTL", "key"))

// 	if ttl != TTL {
// 		t.Fatalf("error, actual: %v, expected: %v", ttl, TTL)
// 	}
// }

// func TestScriptHSETPEX(t *testing.T) {
// 	conn := redisPool.Get()
// 	defer conn.Close()

// 	TTL := 2000.0

// 	script_hsetpex, err := GetScript("1_hsetpex.lua")
// 	if err != nil {
// 		t.Fatalf("error connection to script, %v", err)
// 	}

// 	_, err = script_hsetpex.Do(conn, "key", "field", "value", TTL)
// 	if err != nil {
// 		t.Fatalf("error to hsetpex, %v", err)
// 	}

// 	_ttl, err := redis.Int(conn.Do("PTTL", "key"))
// 	ttl := float64(_ttl)

// 	if ttl > TTL || ttl < TTL*0.99 {
// 		t.Fatalf("error, actual: %v, expected: %v", ttl, TTL)
// 	}
// }

// func TestScriptHINCRBYEX(t *testing.T) {
// 	conn := redisPool.Get()
// 	defer conn.Close()

// 	TTL := 100

// 	conn.Do("HSET", "key", "field", 1)

// 	script_hincrbyex, err := GetScript(1, "hincrbyex.lua")
// 	if err != nil {
// 		t.Fatalf("error connection to script, %v", err)
// 	}

// 	_, err = script_hincrbyex.Do(conn, "key", "field", 99, TTL)
// 	if err != nil {
// 		t.Fatalf("error to hincrbyex, %v", err)
// 	}

// 	ttl, err := redis.Int(conn.Do("TTL", "key"))

// 	if ttl != TTL {
// 		t.Fatalf("error, actual: %v, expected: %v", ttl, TTL)
// 	}

// 	num, err := redis.Int(conn.Do("HGET", "key", "field"))

// 	if num != 100 {
// 		t.Fatalf("error, actual: %v, expected: 100", num)
// 	}
// }

// func TestScriptHINCRBYPEX(t *testing.T) {
// 	conn := redisPool.Get()
// 	defer conn.Close()

// 	TTL := 2000.0

// 	conn.Do("HSET", "key", "field", 1)

// 	script_hincrbyex, err := GetScript(1, "hincrbypex.lua")
// 	if err != nil {
// 		t.Fatalf("error connection to script, %v", err)
// 	}

// 	_, err = script_hincrbyex.Do(conn, "key", "field", 99, TTL)
// 	if err != nil {
// 		t.Fatalf("error to hincrbypex, %v", err)
// 	}

// 	_ttl, err := redis.Int(conn.Do("PTTL", "key"))
// 	ttl := float64(_ttl)

// 	if ttl > TTL || ttl < TTL*0.99 {
// 		t.Fatalf("error, actual: %v, expected: %v", ttl, TTL)
// 	}

// 	num, err := redis.Int(conn.Do("HGET", "key", "field"))

// 	if num != 100 {
// 		t.Fatalf("error, actual: %v, expected: 100", num)
// 	}
// }
