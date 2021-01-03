package rediscript

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const (
	URL = "https://raw.githubusercontent.com/jimako1989/redis-lua-scripts/main/"
)

func GetScript(path string) (*redis.Script, error) {
	splitPath := strings.Split(path, "/")
	keyCount, err := strconv.Atoi(string(splitPath[len(splitPath)-1][0]))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadFile("lua/" + path)
	if err != nil {
		return nil, err
	}

	// res, err := http.Get(URL + path)
	// if err != nil {
	// 	return nil, err
	// }
	// body, err := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// if err != nil {
	// 	return nil, err
	// }

	return redis.NewScript(keyCount, string(body)), nil
}
