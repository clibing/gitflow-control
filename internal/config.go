package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

const CONFIG_YAML = ".control.yaml"

var config Config

type Config struct {
	Issue *Issue `yaml:"issue"` // 用于快速输入jira号
}

type Issue struct {
	FirstEnable bool    `yaml:"first-enable"` // issue号 放在消息的首位
	LeftMarker  string  `yaml:"left-marker"`  // issue号 左边默认标记
	RightMarker string  `yaml:"right-marker"` // issue号 右边默认标记
	Value       []Value `yaml:"value"`        // issue号 或者 jira号
}

type Value struct {
	Number string `yaml:"name"`  // issue 号
	Title  string `yaml:"title"` // issue 描述
}

func init() {
	f := getConfigFilePath()
	exist, err := checkFile(f)
	if !exist || err != nil {
		initDefaultConfig()
		if !exist {
			y, err := yaml.Marshal(&config)
			if err == nil {
				os.WriteFile(f, y, 0644)
			}
		}
		return
	}

	data, err := os.ReadFile(f)
	if err != nil {
		panic(fmt.Errorf("读取配置文件异常, err: %s", err))
	}
	// config.Issue.Value = make([]*Value, 0)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Errorf("加载配置文件异常, err: %s", err))
	}
}

func getConfigFilePath() string {
	homedir, err := homedir.Dir()
	if err != nil {
		panic(fmt.Errorf("获取当前用户的工作目录异常, err: %s", err))
	}
	return filepath.Join(homedir, CONFIG_YAML)
}

func RecoverConfigFile() {
	f := getConfigFilePath()
	exist, _ := checkFile(f)
	if !exist {
		y, err := yaml.Marshal(&config)
		if err == nil {
			os.WriteFile(f, y, 0644)
		}
	}
}

func checkFile(f string) (bool, error) {
	// 检查文件是否存在
	_, err := os.Stat(f)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func initDefaultConfig() {
	v := make([]Value, 0)
	v = append(v, Value{
		Number: "",
		Title:  "",
	})

	config = Config{
		Issue: &Issue{
			FirstEnable: false,
			LeftMarker:  "",
			RightMarker: "",
			Value:       v,
		},
	}

}

func GetConfig() *Config {
	return &config
}
