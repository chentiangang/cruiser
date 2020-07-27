package url

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/chentiangang/cruiser/modules"
)

type TriggerBy struct {
	// 当返回的值包含改字段的字符串 代表接口正常，否则报警
	ResponseContains []string `yaml:"response_contains,omitempty"`

	// 当返回值和改正则表达式能匹配到结果 代表接口正常，否则报警
	MatchRe string `yaml:"match_re,omitempty"`

	// 当返回值等于expected 代表接口正常，否则报警
	Expected string `yaml:"expected,omitempty"`

	// 当该表达式的结果为true 代表接口正常，否则报警
	Expr string `yaml:"expr,omitempty"`
}

func (u UrlConfig) trigger(resp *http.Response) modules.Alert {

	var alertMessage string

	switch {
	case u.TriggerBy.byResponseContains():
		alertMessage = u.triggerByResponseContains(resp)
	case u.TriggerBy.byMatchRe():
		alertMessage = u.triggerByMatchRe(resp)
	case u.TriggerBy.byExpected():
		alertMessage = u.triggerByExpected(resp)
	default:
		alertMessage = fmt.Sprintf("%s\nERROR: \"%s\"", u.Url, "Not found the trigger conditions. Check YAML config")
	}

	if alertMessage != "" {
		u.Alert.Annotations["summary"] = alertMessage
		return u.Alert
	}
	return modules.Alert{}
}

func (u UrlConfig) triggerByResponseContains(resp *http.Response) (msg string) {

	body := respBodyString(resp)
	var str []string
	for _, i := range u.TriggerBy.ResponseContains {
		if !strings.Contains(body, i) {
			str = append(str, i)
		}
	}

	if str != nil {
		u.Alert.Labels["ResponseContains"] = fmt.Sprintf("%s", str)
		return fmt.Sprintf("%s\n返回值:\n%s\n不包含:\n%s", u.Url, body, str)
	}
	return
}

func (u UrlConfig) triggerByMatchRe(resp *http.Response) (msg string) {
	body := respBodyString(resp)

	re := regexp.MustCompile(u.TriggerBy.MatchRe)
	//if len(body) < 200 {
	//	xlog.LogDebug("url: %s, current value: %s", u.Url, body)
	//}
	//xlog.LogDebug("url: %s, match value: %s", u.Url, re.FindString(body))
	if re.FindString(body) == "" {
		msg = fmt.Sprintf("%s\n"+
			"测试用例:\n %s\n"+
			"当前值:\n %s",
			u.Url, u.Request.Data, body)
		return msg
	}
	return
}

func (u UrlConfig) triggerByExpected(resp *http.Response) (msg string) {
	body := respBodyString(resp)
	if body != u.TriggerBy.Expected {
		return fmt.Sprintf("%s\n当前值:\n%s\n,预期值:\n%s\n,statusCode:\n%d", u.Url, body, u.TriggerBy.Expected, resp.Status)
	}
	return
}

func (u UrlConfig) triggerByExpr(resp *http.Response) (msg string) {
	return
}

func (t TriggerBy) byResponseContains() bool {
	return t.ResponseContains != nil
}

func (t TriggerBy) byMatchRe() bool {
	return t.MatchRe != ""
}

func (t TriggerBy) byExpected() bool {
	return t.Expected != ""
}

func respBodyString(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(body)
}
