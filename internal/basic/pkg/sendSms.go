package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io/ioutil"
	"ito-deposit/internal/conf"
	"net/http"
	"net/url"
)

type SendSms struct {
	conf *conf.Data
}

// NewData .
func NewSendSms(c *conf.Data, logger log.Logger) (*SendSms, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &SendSms{
		conf: c,
	}, cleanup, nil
}

func (s *SendSms) SendSms(phone string, code int) bool {
	content := fmt.Sprintf("【短信宝】您的验证码是%d,有效时间为五分钟。", code)
	endata := url.QueryEscape(content)

	urlSms := fmt.Sprintf("https://api.smsbao.com/sms?u=%s&p=%s&m=%s&c=%s",
		"13293991000", s.conf.Smscode, phone, endata)

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", urlSms, nil)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var resultCode int
	err = json.Unmarshal(body, &resultCode)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if resultCode == 0 {
		return true
	}
	fmt.Println(resultCode)
	return false
}
