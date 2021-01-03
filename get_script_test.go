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

func TestScriptTTLAT(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	script_ttlat, err := GetScript("util/1_ttlat.lua")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	_, err = conn.Do("SETEX", "Alice", fmt.Sprint(TTL), "Bob")
	if err != nil {
		t.Fatalf("failed to SETEX, %v", err)
	}

	now := time.Now().Unix()
	expireAt, err := redis.Int64(script_ttlat.Do(conn, "Alice"))
	if err != nil {
		t.Fatalf("error while doing script, %v", err)
	}

	if expireAt != now+TTL {
		t.Fatalf("error, expireAt: %v, time+TTL: %v", expireAt, now+TTL)
	}
}

func TestScriptHSETXP(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("hashes_xp/2_hsetxp.lua")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), "field", "value"))
	if err != nil {
		t.Fatalf("error to hsetxp, %v", err)
	}
	if !r {
		t.Fatalf("failed hsetxp")
	}

	b, err := redis.Bool(conn.Do("HEXISTS", "key", "field"))
	if !b {
		t.Fatalf("failed to HEXISTS because the field doesn't exist, should exist")
	}
}

func TestScriptHMGETXP(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	TestScriptHSETXP(t)

	script, err := GetScript("hashes_xp/1_hmgetxp.lua")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	s, err := redis.Strings(script.Do(conn, "key", "field"))
	if err != nil {
		t.Fatalf("error to hmgetxp, %v", err)
	}

	if s[0] != "value" {
		t.Fatalf("can't find the field. expect: value, but actual: %v", s)
	}

	time.Sleep(time.Duration(TTL+1) * time.Second)

	s, err = redis.Strings(script.Do(conn, "key", "field"))
	if err != nil {
		t.Fatalf("error to hmgetxp, %v", err)
	}

	if len(s) != 0 {
		t.Fatalf("found the value. expect: empty, but actual: %v", s)
	}
}
