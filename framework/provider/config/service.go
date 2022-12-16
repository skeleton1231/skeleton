package config

import (
	"bytes"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

// 删除文件的操作
func (conf *HadeConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func NewHadeConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	// check file exist or not
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}

	// Initialization
	hadeConf := &HadeConfig{
		c:        container,
		folder:   envFolder,
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		keyDelim: ".",
		lock:     sync.Mutex{},
	}

	// Read File
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, file := range files {
		fileName := file.Name()
		err := hadeConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// Monitor Folder File
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						hadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						hadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						hadeConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	return hadeConf, nil
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

// Get config
func (conf *HadeConfig) Get(key string) interface{} {
	return conf.find(key)
}

// GetBool get bool config
func (conf *HadeConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

// GetInt get int
func (conf *HadeConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *HadeConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *HadeConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *HadeConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *HadeConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *HadeConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *HadeConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *HadeConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *HadeConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be an pointer
func (conf *HadeConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(conf.find(key))
}
