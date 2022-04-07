package request

import (
	"crypto/tls"
	"encoding/json"
	"github.com/fusuwei/gspider/pkg/response"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Session struct {
	Host           string
	Timeout        time.Duration
	Verify         bool // 是否验证证书
	AllowRedirects bool // 禁止跳转
	UA             string
	Header         map[string]string
	cookie         *cookiejar.Jar
	Client         *http.Client
}

func New(host string, timeout time.Duration, verify bool, allowRedirects bool, header map[string]string) *Session {
	session := &Session{
		Host:           host,
		Timeout:        timeout,
		Verify:         verify,
		AllowRedirects: allowRedirects,
		Header:         header,
		cookie:         nil,
	}
	client := &http.Client{
		Timeout: timeout,
	}
	if !allowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: verify,
		},
	}

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	if jar, err := cookiejar.New(&options); err != nil {

	} else {
		client.Jar = jar
	}
	session.Client = client
	return session
}

func (s *Session) Request(req *Request) (*response.Response, error) {
	var (
		body        io.Reader
		contentType string
	)
	u := req.Url.String()
	method := strings.ToUpper(req.Method)
	if req.Data != nil {
		data, err := json.Marshal(req.Data)
		if err != nil {
			return nil, err
		}
		contentType = "application/x-www-form-urlencoded"
		body = strings.NewReader(string(data))
	} else if req.Json != nil {
		data, err := json.Marshal(req.Json)
		if err != nil {
			return nil, err
		}
		contentType = "application/json;charset=UTF-8"
		body = strings.NewReader(string(data))
	}
	request, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("content-type", contentType)
	for k, v := range s.Header {
		if strings.ToLower(k) == "user-agent" && s.UA != "" {
			request.Header.Set("user-agent", s.UA)
			continue
		}
		request.Header.Set(k, v)
	}
	if request.Header.Get("user-agent") == "" && s.UA != "" {
		request.Header.Set("user-agent", s.UA)
	}

	res, err := s.Client.Do(request)
	if err != nil {
		return nil, err
	}
	resp := response.New(res, req.Callback, req.Meta)
	return resp, nil
}
