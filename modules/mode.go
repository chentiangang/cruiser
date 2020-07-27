package modules

import (
	"time"

	"github.com/prometheus/common/model"
)

type Config struct {
	Global    Global    `yaml:"global"`
	UrlConfig UrlConfig `yaml:"url_config"`
	Task      chan Cruiser
}

type UrlConfig struct {
	Include []string `yaml:"include"`
}

type Global struct {
	AlertmanagerApi string         `yaml:"alertmanager_api"`
	Interval        model.Duration `yaml:"interval,omitempty"`
	Timeout         model.Duration `yaml:"timeout,omitempty"`
}

type Alert struct {
	Labels       map[string]string `yaml:"labels,omitempty" json:"labels"`
	Annotations  map[string]string `yaml:"annotations,omitempty" json:"annotations"`
	GeneratorURL string            `yaml:"generatorURL" json:"generatorURL"`
	//StartsAt     time.Time         `json:"startsAt"`
	//EndsAt       time.Time         `json:"endsAt"`
	//isEnd chan bool
}

type Cruiser interface {
	// 加载配置文件
	LoadConfig(cfg *Config) error

	// 获取请求间隔
	GetInterval() time.Duration

	// 发送请求，返回报警信息
	SendRequest() Alert
}
