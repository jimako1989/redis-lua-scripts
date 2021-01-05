package rediscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestHSETXP(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("HASHES_XP/2_HSETXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), "field", "value"))
	if err != nil {
		t.Fatalf("error to HSETXP, %v", err)
	}
	if !r {
		t.Fatalf("failed HSETXP")
	}

	b, err := redis.Bool(conn.Do("HEXISTS", "key", "field"))
	if !b {
		t.Fatalf("failed to HEXISTS because the field doesn't exist, should exist")
	}
}

func TestHMGETXP(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	TestHSETXP(t)

	script, err := GetScript("HASHES_XP/1_HMGETXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	s, err := redis.Strings(script.Do(conn, "key", "field"))
	if err != nil {
		t.Fatalf("error to HMGETXP, %v", err)
	}

	if s[0] != "value" {
		t.Fatalf("can't find the field. expect: value, but actual: %v", s)
	}

	time.Sleep(time.Duration(TTL+1) * time.Second)

	s, err = redis.Strings(script.Do(conn, "key", "field"))
	if err != nil {
		t.Fatalf("error to HMGETXP, %v", err)
	}

	if len(s) != 0 {
		t.Fatalf("found the value. expect: empty, but actual: %v", s)
	}
}
