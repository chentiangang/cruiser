package url

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/chentiangang/cruiser/modules"

	"github.com/chentiangang/xlog"
	"github.com/prometheus/common/model"
)

type UrlConfig struct {
	Url     string `yaml:"url"`
	Request struct {
		Method   string            `yaml:"method"`
		Data     string            `yaml:"data"`
		Header   map[string]string `yaml:"add_header,omitempty"`
		Interval model.Duration    `yaml:"interval,omitempty"`
		Timeout  model.Duration    `yaml:"timeout,omitempty"`
	} `yaml:"request"`
	// 这个结构包含所有触发条件的实现
	TriggerBy TriggerBy `yaml:"trigger_by"`

	Alert modules.Alert `yaml:"alert"`
}

func (u UrlConfig) newHttpClient() *http.Client {
	t, err := time.ParseDuration(u.Request.Timeout.String())
	if err != nil {
		panic(err)
	}
	return &http.Client{
		Timeout: t,
	}
}

func (u UrlConfig) addHeader(req *http.Request) {
	for key, value := range u.Request.Header {
		req.Header.Add(key, value)
	}
}

func (u UrlConfig) do() (resp *http.Response, err error) {
	client := u.newHttpClient()

	req, _ := http.NewRequest(u.Request.Method, u.Url, bytes.NewReader([]byte(u.Request.Data)))
	u.addHeader(req)
	return client.Do(req)
}

func (u UrlConfig) SendRequest() modules.Alert {
	resp, err := u.do()
	xlog.LogDebug("requst url: %s", u.Url)
	if err != nil {
		u.Alert.Annotations["summary"] = fmt.Sprintf("URL: %s,\nERROR: \"%s\"", u.Url, err)
		return u.Alert
	}

	return u.trigger(resp)
}
