package rediscript

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func GetScript(path string) (*redis.Script, error) {
	splitPath := strings.Split(path, "/")
	keyCount, err := strconv.Atoi(string(splitPath[len(splitPath)-1][0]))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadFile("lua/" + strings.ToUpper(path) + ".lua")
	if err != nil {
		return nil, err
	}

	return redis.NewScript(keyCount, string(body)), nil
}

func GetAllScripts(group string) (map[string]*redis.Script, error) {
	files, err := ioutil.ReadDir("lua/" + string(group))
	if err != nil {
		return nil, err
	}

	scripts := make(map[string]*redis.Script)
	for _, file := range files {
		name := file.Name()
		script, err := GetScript(string(group) + "/" + name[:len(name)-4])
		if err != nil {
			return nil, err
		}
		scripts[name[2:len(name)-4]] = script
	}
	return scripts, nil
}
