package rediscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestGetAllScripts(t *testing.T) {
	scripts, err := GetAllScripts("SORTED_SETS_XP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}
	if len(scripts) != 17 {
		t.Fatalf("invalid the number of scripts, expect: 13, actual: %v", len(scripts))
	}
	fmt.Print(scripts["ZADDXP"])
}

func TestTTLAT(t *testing.T) {
	script_ttlat, err := GetScript("UTIL/1_TTLAT")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	conn := redisPool.Get()
	defer conn.Close()

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
