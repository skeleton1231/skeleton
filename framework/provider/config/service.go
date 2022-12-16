package config

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skeleton1231/skeleton/framework"
	"github.com/skeleton1231/skeleton/framework/contract"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

// HadeConfig
type HadeConfig struct {
	c        framework.Container
	folder   string
	keyDelim string                 // path delimeter, default:"."
	lock     sync.Mutex             // config read&write lock
	envMaps  map[string]string      //all env args
	confMaps map[string]interface{} // config file struct, key as file name
	confRaws map[string][]byte      // config file raw data
}

// read specific config file
func (conf *HadeConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	// check yml or yaml file
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		// read file content
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		// replace the env args in text
		bf = replace(bf, conf.envMaps)
		// unmarshal file
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		// rear app.path data and update the relevant file
		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := conf.c.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

// replace: use envMaps to replace context env(xxx) args
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	// use replaceAll as simple solution
	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}

// search specific path config
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	// iterate
	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}

		// next
		switch next.(type) {
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}
	return nil
}

// get config item through path
func (conf *HadeConfig) find(key string) interface{} {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}
