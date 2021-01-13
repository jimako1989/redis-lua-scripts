package rediscript

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func setup_SortedListsXP() {
	conn := redisPool.Get()
	defer conn.Close()
	script, _ := GetScript("SORTED_LISTS_XP/1_ZLDELXP")
	b, _ := redis.Bool(script.Do(conn, "keyy", 0, -1))
	if !b {
		panic("Failed to ZLDELXP")
	}
}

func TestZLPUSHXP(t *testing.T) {
	setup_SortedListsXP()

	script, err := GetScript("SORTED_LISTS_XP/2_ZLPUSHXP")
	if err != nil {
		t.Fatalf("error connection to script, %v", err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	r, err := redis.Bool(script.Do(conn, "keyy", fmt.Sprint(TTL), 10.0, "fieldd"))
	if err != nil {
		t.Fatalf("error to ZLPUSHXP, %v", err)
	}
	if !r {
		t.Fatalf("failed ZLPUSHXP")
	}

	script, _ = GetScript("SORTED_LISTS_XP/3_ZLRANGEXP")
	s, err := redis.Strings(script.Do(conn, "keyy", 0, -1))
	if err != nil {
		t.Fatalf("error to ZLRANGEXP, %v", err)
	}
	if len(s) != 1 {
		t.Fatalf("expect: 1, gots %v", len(s))
	}
	if s[0] != "fieldd" {
		t.Fatalf("expect: fieldd, gots %v", s)
	}
}

func TestZLPOPXP(t *testing.T) {
	TestZLPUSHXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/1_ZLPOPXP")
	s, err := redis.String(script.Do(conn, "keyy"))
	if err != nil {
		t.Fatalf("error to ZLPOPXP, %v", err)
	}
	if s != "fieldd" {
		t.Fatalf("expect: fieldd, gots %v", s)
	}
}

func TestExpire(t *testing.T) {
	TestZLPUSHXP(t)

	time.Sleep(time.Duration(TTL+1) * time.Second)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/3_ZLRANGEXP")
	s, err := redis.Strings(script.Do(conn, "keyy", 0, -1))
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
	_, err := redis.Bool(script.Do(conn, "keyy", fmt.Sprint(TTL), 123, "fieldd2"))
	if err != nil {
		t.Fatalf("error to ZLPOPXP, %v", err)
	}

	script, _ = GetScript("SORTED_LISTS_XP/3_ZLRANGEBYSCOREXP")
	s, err := redis.Strings(script.Do(conn, "keyy", 0, 200))
	if err != nil {
		t.Fatalf("error to ZLRANGEXP, %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("expect: 2, gots %v. %v", len(s), s)
	}
	if s[0] != "fieldd" && s[1] != "fieldd2" {
		t.Fatalf("expect: [fieldd, fieldd2], gots %v", s)
	}
}

func TestZLREVRANGEBYSCOREXP(t *testing.T) {
	TestZLPUSHXP(t)

	conn := redisPool.Get()
	defer conn.Close()

	script, _ := GetScript("SORTED_LISTS_XP/2_ZLPUSHXP")
	_, err := redis.Bool(script.Do(conn, "keyy", fmt.Sprint(TTL), 123, "fieldd2"))
	if err != nil {
		t.Fatalf("error to ZLPOPXP, %v", err)
	}

	script, _ = GetScript("SORTED_LISTS_XP/3_ZLREVRANGEBYSCOREXP")
	s, err := redis.Strings(script.Do(conn, "keyy", 200, 0))
	if err != nil {
		t.Fatalf("error to ZLREVRANGEBYSCOREXP, %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("expect: 2, gots %v. %v", len(s), s)
	}
	if s[0] != "fieldd2" && s[1] != "fieldd" {
		t.Fatalf("expect: [fieldd2, fieldd], gots %v", s)
	}
}
