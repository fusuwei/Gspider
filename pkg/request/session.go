package request

import (
	"crypto/tls"
	"encoding/json"
	"github.com/fusuwei/gspider/pkg/response"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Session struct {
	Timeout        time.Duration
	Verify         bool // 是否验证证书
	AllowRedirects bool // 是否禁止跳转
	Client         *http.Client
	Transport      *http.Transport
}

func New(timeout time.Duration, verify bool, allowRedirects bool) *Session {
	session := &Session{
		Timeout:        timeout,
		Verify:         verify,
		AllowRedirects: allowRedirects,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: verify,
			},
		},
	}
	client := &http.Client{
		Timeout: timeout,
	}
	if !allowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
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
	// 设置headers
	request.Header.Set("content-type", contentType)
	for k, v := range req.Header {
		if strings.ToLower(k) == "user-agent" && req.UA != "" {
			request.Header.Set("user-agent", req.UA)
			continue
		}
		request.Header.Set(k, v)
	}
	if request.Header.Get("user-agent") == "" && req.UA != "" {
		request.Header.Set("user-agent", req.UA)
	}

	transport := s.Transport.Clone()
	if req.Proxy != "" {
		if p, err := url.Parse(req.Proxy); err == nil {
			transport.Proxy = http.ProxyURL(p)
		}
	}
	s.Client.Transport = transport
	res, err := s.Client.Do(request)
	if err != nil {
		return nil, err
	}
	resp := response.New(req.Url, res, req.Callback, req.Meta)
	return resp, nil
}
