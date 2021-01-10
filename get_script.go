package rediscript

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func GetScript(filepath string) (*redis.Script, error) {
	splitPath := strings.Split(filepath, "/")
	keyCount, err := strconv.Atoi(string(splitPath[len(splitPath)-1][0]))
	if err != nil {
		return nil, err
	}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to call runtime.Caller(0)")
	}
	body, err := ioutil.ReadFile(path.Dir(file) + "/lua/" + strings.ToUpper(filepath) + ".lua")
	if err != nil {
		return nil, err
	}

	return redis.NewScript(keyCount, string(body)), nil
}

func GetAllScripts(group string) (map[string]*redis.Script, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to call runtime.Caller(0)")
	}
	files, err := ioutil.ReadDir(path.Dir(file) + "/lua/" + string(group))
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
