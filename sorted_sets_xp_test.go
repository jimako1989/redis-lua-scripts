package rediscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	score123 = 123.0
	score234 = 234.5
)

func TestZADDXP(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/2_ZADDXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), score123, "member"))
	if err != nil {
		t.Fatalf("error to ZADDXP, %v", err)
	}
	if !r {
		t.Fatalf("failed ZADDXP")
	}

	f, err := redis.Float64(conn.Do("ZSCORE", "key", "member"))
	if err != nil {
		t.Fatalf("failed to ZSCORE because the member doesn't exist, should exist")
	}
	if f != score123 {
		t.Fatalf("failed to ZSCORE. expect score: %v, but actual: %v", score123, f)
	}
}

func TestZCARDXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/1_ZCARDXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	c, err := redis.Int64(script.Do(conn, "key"))
	if err != nil {
		t.Fatalf("error to ZCARDXP, %v", err)
	}
	if c != 1 {
		t.Fatalf("failed ZCARDXP. expect: 1 but actual: %v", c)
	}

	time.Sleep(time.Duration(TTL+1) * time.Second)

	c, err = redis.Int64(script.Do(conn, "key"))
	if err != nil {
		t.Fatalf("error to ZCARDXP, %v", err)
	}
	if c != 0 {
		t.Fatalf("failed ZCARDXP. expect: 0 but actual: %v", c)
	}
}

func TestZCOUNTXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/3_ZCOUNTXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	c, err := redis.Int64(script.Do(conn, "key", 0, -1))
	if err != nil {
		t.Fatalf("error to ZCOUNTXP, %v", err)
	}
	if c != 1 {
		t.Fatalf("failed ZCOUNTXP. expect: 1 but actual: %v", c)
	}

	time.Sleep(time.Duration(TTL+1) * time.Second)

	c, err = redis.Int64(script.Do(conn, "key", 0, -1))
	if err != nil {
		t.Fatalf("error to ZCOUNTXP, %v", err)
	}
	if c != 0 {
		t.Fatalf("failed ZCOUNTXP. expect: 0 but actual: %v", c)
	}
}

func TestZSCOREXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/2_ZSCOREXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	score, err := redis.Float64(script.Do(conn, "key", "member"))
	if err != nil {
		t.Fatalf("error to ZSCOREXP, %v", err)
	}
	if score != score123 {
		t.Fatalf("failed ZSCOREXP. expect: %v but actual: %v", score123, score)
	}

	time.Sleep(time.Duration(TTL+1) * time.Second)

	score, err = redis.Float64(script.Do(conn, "key", "member"))
	if err.Error() != "redigo: nil returned" {
		t.Fatalf("error to ZSCOREXP, %v", err)
	}
}

func TestZEXPIREATXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/2_ZEXPIREATXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	now := time.Now().Unix()
	expireAt, err := redis.Int64(script.Do(conn, "key", "member"))
	if err != nil {
		t.Fatalf("error to ZEXPIREATXP, %v", err)
	}
	if expireAt != now+TTL {
		t.Fatalf("error, expireAt: %v, time+TTL: %v", expireAt, now+TTL)
	}
}

func TestZRANGEXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/2_ZADDXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), score234, "member2"))
	if err != nil {
		t.Fatalf("error to ZADDXP, %v", err)
	}
	if !r {
		t.Fatalf("failed ZADDXP")
	}

	script, err = GetScript("SORTED_SETS_XP/3_ZRANGEXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	s, err := redis.Strings(script.Do(conn, "key", 0, -1, "WITHSCORES"))
	if err != nil {
		t.Fatalf("failed to ZRANGEXP, %v", err)
	}
	if len(s) != 4 {
		t.Fatalf("error to ZRANGEXP, expect #member: 4, but actual: %v", s)
	}
	if s[1] != "123" {
		t.Fatalf("error to ZRANGEXP, expect %v: %v, but actual: %v", "member", "123", s[1])
	}
	if s[3] != "234.5" {
		t.Fatalf("error to ZRANGEXP, expect %v: %v, but actual: %v", "member2", "234.5", s[3])
	}
}

func TestZREVRANGEXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/2_ZADDXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), score234, "member2"))
	if err != nil {
		t.Fatalf("error to ZADDXP, %v", err)
	}
	if !r {
		t.Fatalf("failed ZADDXP")
	}

	script, err = GetScript("SORTED_SETS_XP/3_ZREVRANGEXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	s, err := redis.Strings(script.Do(conn, "key", 0, -1, "WITHSCORES"))
	if err != nil {
		t.Fatalf("failed to ZRANGEXP, %v", err)
	}
	if len(s) != 4 {
		t.Fatalf("error to ZRANGEXP, expect #member: 4, but actual: %v", s)
	}
	if s[3] != "123" {
		t.Fatalf("error to ZRANGEXP, expect %v: %v, but actual: %v", "member", "123", s[3])
	}
	if s[1] != "234.5" {
		t.Fatalf("error to ZRANGEXP, expect %v: %v, but actual: %v", "member2", "234.5", s[1])
	}
}

func TestZRANGEBYSCOREXP(t *testing.T) {
	TestZADDXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, err := GetScript("SORTED_SETS_XP/2_ZADDXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	r, err := redis.Bool(script.Do(conn, "key", fmt.Sprint(TTL), score234, "member2"))
	if err != nil {
		t.Fatalf("error to ZADDXP, %v", err)
	}
	if !r {
		t.Fatalf("failed ZADDXP")
	}

	script, err = GetScript("SORTED_SETS_XP/3_ZRANGEBYSCOREXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	s, err := redis.Strings(script.Do(conn, "key", 122.9, 234.6))
	if err != nil {
		t.Fatalf("failed to ZRANGEBYSCOREXP, %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("error to ZRANGEBYSCOREXP, expect #member: 2, but actual: %v", s)
	}
	if s[0] != "member" {
		t.Fatalf("error to ZRANGEBYSCOREXP, expect: %v, but actual: %v", "member", s[0])
	}
	if s[1] != "member2" {
		t.Fatalf("error to ZRANGEBYSCOREXP, expect: %v, but actual: %v", "member2", s[1])
	}
}
