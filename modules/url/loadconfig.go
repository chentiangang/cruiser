package url

import (
	"io/ioutil"
	"time"

	"github.com/chentiangang/cruiser/modules"

	yaml "gopkg.in/yaml.v2"
)

//func loadConfig() {
//	var cfg
//}

func (u UrlConfig) LoadConfig(cfg *modules.Config) (err error) {
	for _, i := range cfg.UrlConfig.Include {
		bs, err := ioutil.ReadFile(i)
		if err != nil {
			panic(err)
		}

		var urls []UrlConfig
		err = yaml.Unmarshal(bs, &urls)
		if err != nil {
			panic(err)
		}

		for _, i := range urls {
			i.setInterval(cfg.Global.Interval.String())
			i.setTimeout(cfg.Global.Timeout.String())
			i.setDefaultAlert()
			//xlog.LogDebug("%+v", i.Url)
			cfg.Task <- i
		}
	}

	return nil
}

func (u *UrlConfig) setInterval(interval string) {
	if u.Request.Interval.String() == "0s" {
		if interval == "" {
			u.Request.Interval.Set("1m")
		}
		u.Request.Interval.Set(interval)
	}
}

func (u *UrlConfig) setTimeout(timeout string) {
	if u.Request.Timeout.String() == "0s" {
		if timeout == "0s" {
			u.Request.Timeout.Set("10s")
		}
		u.Request.Timeout.Set(timeout)
	}
}

func (u *UrlConfig) setDefaultAlert() {
	u.Alert.Labels["url"] = u.Url
	u.Alert.GeneratorURL = "http://localhost:8080/metrics"

	if len(u.Alert.Annotations) == 0 {
		u.Alert.Annotations = make(map[string]string, 2)
	}
}

func (u UrlConfig) GetInterval() time.Duration {
	duration, err := time.ParseDuration(u.Request.Interval.String())
	if err != nil {
		panic(err)
	}
	return duration
}
