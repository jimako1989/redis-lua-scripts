package rediscript

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest/v3"
)

var (
	dockerRes *dockertest.Resource
	conn      redis.Conn
)

func TestMain(m *testing.M) {
	setup()

	defer dockerRes.Close()
	os.Exit(m.Run())
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

	for {
		conn, err = redis.DialURL(fmt.Sprintf("redis://localhost:%s", dockerRes.GetPort("6379/tcp")))
		if err == nil {
			break
		} else {
			time.Sleep(time.Second)
		}
	}
}

func TestScriptTTLAT(t *testing.T) {
	TTL := int64(100)
	defer conn.Close()

	script_ttlat, err := GetScript(1, "ttlat.lua")
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
