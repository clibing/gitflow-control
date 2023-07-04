package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

var config *Config

type Config struct {
	Mode   string    `yaml:"mode"`   // issue号 放在消息的首位, 模式， 默认为auto:根据PrefixUrl自动调整, first: 放在首位, standard: 按照规范设置
	Issue  *Issue    `yaml:"issue"`  // 用于快速输入jira号
	Record []*Record `yaml:"record"` // 记录器
}

// 记录器
type Record struct {
	Project  string            `yaml:"project"`  // 对应的Git的项目名字
	Describe map[string]string `yaml:"describe"` // 分支的描述信息
}

type Issue struct {
	PrefixUrl   []string                           `yaml:"prefix-url"`      // 如果配置默认为自动匹配模式
	LeftMarker  string                             `yaml:"left-marker"`     // issue号 左边默认标记
	RightMarker string                             `yaml:"right-marker"`    // issue号 右边默认标记
	Value       string                             `yaml:"value,omitempty"` // issue号 或者 jira号
	List        []*RecordIssue                     `yaml:"list"`            // 最近的issue号
	Describes   map[string]map[string]*BugDescribe `yaml:"listv2"`          // 最近的issue号
}

type RecordIssue struct {
	ProjectName string `yaml:"project-name"` // 对应的Git的项目名字
	Number      string `yaml:"number"`       // 对应的issue号
}

// bug 描述
type BugDescribe struct {
	Number   string `yaml:"number"`   // jira number
	Describe string `yaml:"describe"` // jira description
	Time     string `yaml:"time"`     // create time
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
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(fmt.Errorf("加载配置文件异常, err: %s", err))
	}
	config = &c
}

func getConfigFilePath() string {
	homedir, err := homedir.Dir()
	if err != nil {
		panic(fmt.Errorf("获取当前用户的工作目录异常, err: %s", err))
	}
	return filepath.Join(homedir, HomeDir, ConfigYaml)
}

func RecoverConfigFile() {
	f := getConfigFilePath()
	exist, _ := checkFile(f)
	if !exist {
		y, err := yaml.Marshal(config)
		if err == nil {
			err = os.WriteFile(f, y, 0644)
			if err != nil {
				fmt.Println(err.Error())
			}
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
	config = &Config{
		Mode: "auto",
		Issue: &Issue{
			LeftMarker:  "",
			RightMarker: "",
			Value:       "",
		},
	}

}

func GetConfig() *Config {
	return config
}

func RequiredFooter() bool {
	switch config.Mode {
	case "first":
		return true
	case "standard":
		return false
	// case "auto":
	case "auto":
		url, err := GetOriginUrl()
		if err != nil {
			return false
		}
		for _, prefix := range config.Issue.PrefixUrl {
			if strings.HasPrefix(url, prefix) {
				return true
			}
		}
		return false
	}
	return false
}

func Rewrite() {
	f := getConfigFilePath()
	y, err := yaml.Marshal(config)
	if err == nil {
		err = os.WriteFile(f, y, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
