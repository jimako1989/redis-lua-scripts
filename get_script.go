package rediscript

import (
	"io/ioutil"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

const (
	URL = "https://raw.githubusercontent.com/jimako1989/redis-lua-scripts/main/"
)

func GetScript(keyCount int, fname string) (*redis.Script, error) {
	res, err := http.Get(URL + fname)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return redis.NewScript(keyCount, string(body)), nil
}
