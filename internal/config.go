package internal

type Config struct {
	Issue *Issue `yaml:"issue"` // 用于快速输入jira号
}

type Issue struct {
	FirstEnable bool     `yaml:"first-enable"` // issue号 放在消息的首位
	LeftMarker  string   `yaml:"left-marker"`  // issue号 左边默认标记
	RightMarker string   `yaml:"right-marker"` // issue号 右边默认标记
	Value       []string `yaml:"value"`        // issue号 或者 jira号
}

