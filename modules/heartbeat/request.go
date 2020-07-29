package heartbeat

import (
	"cruiser/modules"
	"time"

	"github.com/prometheus/common/model"
)

type HeartBConfig struct {
	Interval model.Duration
}

func (h HeartBConfig) LoadConfig(cfg *modules.Config) error {
	return nil
}

func (h HeartBConfig) SendRequest() modules.Alert {
	return modules.Alert{}

}

func (h HeartBConfig) GetInterval() time.Duration {
	duration, err := time.ParseDuration(h.Interval.String())
	if err != nil {
		panic(err)
	}
	return duration
}
