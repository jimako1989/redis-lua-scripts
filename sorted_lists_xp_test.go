package rediscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestZLPUSHXP(t *testing.T) {
	script, err := GetScript("SORTED_LISTS_XP/2_ZLPUSHXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), 10.0, "field"))
	if err != nil {
		t.Fatalf("error to ZLPUSHXP, %v", err)
	}
	if !r {
		t.Fatalf("failed ZLPUSHXP")
	}

	script, _ = GetScript("SORTED_LISTS_XP/3_ZLRANGEXP")
	s, err := redis.Strings(script.Do(conn, "key", 0, -1))
	if err != nil {
		t.Fatalf("error to ZLRANGEXP, %v", err)
	}
	if s[0] != "field" {
		t.Fatalf("expect: field, gots %v", s)
	}
}

func TestZLPOPXP(t *testing.T) {
	TestZLPUSHXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/1_ZLPOPXP")
	s, err := redis.String(script.Do(conn, "key"))
	if err != nil {
		t.Fatalf("error to ZLPOPXP, %v", err)
	}
	if s != "field" {
		t.Fatalf("expect: field, gots %v", s)
	}
}

func TestExpire(t *testing.T) {
	TestZLPUSHXP(t)

	time.Sleep(time.Duration(TTL+1) * time.Second)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/3_ZLRANGEXP")
	s, err := redis.Strings(script.Do(conn, "key", 0, -1))
	if err != nil {
		t.Fatalf("error to ZLRANGEXP, %v", err)
	}
	if len(s) != 0 {
		t.Fatalf("expect: zero list, gots %v", s)
	}
}

func TestZLRANGEBYSCOREXP(t *testing.T) {
	TestZLPUSHXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/2_ZLPUSHXP")
	_, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), 123, "field2"))
	if err != nil {
		t.Fatalf("error to ZLPOPXP, %v", err)
	}

	script, _ = GetScript("SORTED_LISTS_XP/3_ZLRANGEBYSCOREXP")
	s, err := redis.Strings(script.Do(conn, "key", 0, 200))
	if err != nil {
		t.Fatalf("error to ZLRANGEXP, %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("expect: 2, gots %v. %v", len(s), s)
	}
	if s[0] != "field" && s[1] != "field2" {
		t.Fatalf("expect: [field, field2], gots %v", s)
	}
}

func TestZLREVRANGEBYSCOREXP(t *testing.T) {
	TestZLPUSHXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/2_ZLPUSHXP")
	_, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), 123, "field2"))
	if err != nil {
		t.Fatalf("error to ZLPOPXP, %v", err)
	}

	script, _ = GetScript("SORTED_LISTS_XP/3_ZLREVRANGEBYSCOREXP")
	s, err := redis.Strings(script.Do(conn, "key", 200, 0))
	if err != nil {
		t.Fatalf("error to ZLREVRANGEBYSCOREXP, %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("expect: 2, gots %v. %v", len(s), s)
	}
	if s[0] != "field2" && s[1] != "field" {
		t.Fatalf("expect: [field2, field], gots %v", s)
	}
}
