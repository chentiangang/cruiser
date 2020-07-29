package url

import (
	"bytes"
	"cruiser/modules"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/chentiangang/xlog"
	"github.com/prometheus/common/model"
)

type UrlConfig struct {
	Url     string `yaml:"url"`
	Request struct {
		Method         string            `yaml:"method"`
		Data           string            `yaml:"data"`
		Header         map[string]string `yaml:"add_header,omitempty"`
		ReplaceOrderId bool              `yaml:"replace_order_id,omitempty"`
		Signature      bool              `yaml:"signature,omitempty"`
		ApiMethod      string            `yaml:"api_method,omitempty"`
		Interval       model.Duration    `yaml:"interval,omitempty"`
		Timeout        model.Duration    `yaml:"timeout,omitempty"`
		AccountId      string            `yaml:"account_id"`
		SecretKeyId    string            `yaml:"secret_key_id"`
		SecretKey      string            `yaml:"secret_key"`
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

func (u *UrlConfig) replaceOrderId() {
	// 生成13位时间戳
	t := time.Now().Add(1000*time.Millisecond).UnixNano() / 1e6
	// 截取accountid
	accountId := strings.Split(u.AccountID(), "-")[1]

	// 随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000)

	// 拼接orderid
	orderId := fmt.Sprintf("%s%d%03d", accountId, t, r)
	repl := fmt.Sprintf("\"orderId\": \"%s\"", orderId)

	//xlog.LogDebug("accountid: %s, time: %d, rand: %d, orderId: %s", accountId, t, r, orderId)
	// 替换
	re := regexp.MustCompile(`"orderId": ".*[0-9]+"`)
	dataByte := re.ReplaceAll([]byte(u.Request.Data), []byte(repl))

	u.Request.Data = string(dataByte)
	//xlog.LogDebug("replace orderId: %s", u.Request.Data)
}

func (u UrlConfig) do() (resp *http.Response, err error) {
	client := u.newHttpClient()

	if u.Request.ReplaceOrderId == true {
		u.replaceOrderId()
	}
	if u.Request.Signature == true {
		u.Url = u.Url + "?signature=" + u.Sig()
	}

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
