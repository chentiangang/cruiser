package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/chentiangang/cruiser/modules"
	"github.com/chentiangang/cruiser/modules/url"

	"github.com/chentiangang/xlog"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	yaml "gopkg.in/yaml.v2"
)

var cfg modules.Config

//const collectTime = time.Second * 10
const filename = "./cruiser.yml"

func NewCruiser(name string) modules.Cruiser {
	if name == "url" {
		return url.UrlConfig{}
	}

	panic(fmt.Sprintf("Not found %s\n", name))
	return nil
}

func main() {
	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.Parse()

	cfg.Task = make(chan modules.Cruiser, 50)

	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bs, &cfg)
	if err != nil {
		panic(err)
	}

	xlog.LogDebug("%+v", cfg)
	go Run()

	if !reflect.DeepEqual(cfg.UrlConfig, modules.UrlConfig{}) {
		cruiser := NewCruiser("url")
		err = cruiser.LoadConfig(&cfg)
		if err != nil {
			panic(err)
		}
	}

	close(cfg.Task)

	http.Handle("/metrics", promhttp.Handler())
	xlog.LogFatal("%s", http.ListenAndServe(*addr, nil))
}

func SendAlert(alert modules.Alert) {
	var Alerts []modules.Alert
	Alerts = append(Alerts, alert)

	bs, err := json.Marshal(Alerts)
	if err != nil {
		xlog.LogError("%s", err)
		panic(err)
	}

	//xlog.LogInfo("%s", bs)
	client := &http.Client{}
	req, err := http.NewRequest("POST", cfg.Global.AlertmanagerApi, bytes.NewReader(bs))

	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	_, err = client.Do(req)
	if err != nil {
		xlog.LogError("alert data: %s, error: %s", bs, err)
		panic(err)
	}
	xlog.LogWarn("send alert: %s", bs)
}

func Run() {
	for task := range cfg.Task {
		xlog.LogDebug("%+v", task)
		go func(i modules.Cruiser) {
			ticker := time.NewTicker(i.GetInterval())
			defer ticker.Stop()

			// 启动完成后即发送一个请求，避免调试等待时间
			alert := i.SendRequest()
			if !reflect.DeepEqual(alert, modules.Alert{}) {
				SendAlert(alert)
			}

			// 定时器
			for {
				select {
				case <-ticker.C:
					alert = i.SendRequest()
					//xlog.LogDebug("%+v", alert)
					if !reflect.DeepEqual(alert, modules.Alert{}) {
						go SendAlert(alert)
					}
				}
			}
		}(task)
	}
}
