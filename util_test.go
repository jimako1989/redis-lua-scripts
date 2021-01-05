package rediscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestTTLAT(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	script_ttlat, err := GetScript("UTIL/1_TTLAT.lua")
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
