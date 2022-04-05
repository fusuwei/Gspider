package request

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Request struct {
	Url      *url.URL
	Method   string
	Data     map[string]interface{}
	Json     map[string]interface{}
	Callback string
	Meta     map[string]interface{}
	Retry    int
}

func (r *Request) ToPublish() string {
	request := make(map[string]string)
	request["url"] = r.Url.String()
	request["method"] = r.Method
	if r.Data != nil {
		if data, err := json.Marshal(r.Data); err != nil {
			request["data"] = ""
		} else {
			request["data"] = string(data)
		}
	} else if r.Json != nil {
		if j, err := json.Marshal(r.Json); err != nil {
			request["data"] = ""
		} else {
			request["data"] = string(j)
		}
	}

	request["retry"] = strconv.Itoa(r.Retry)
	request["callback"] = r.Callback
	data, err := json.Marshal(request)
	if err != nil {
		return ""
	}
	return string(data)
}

func ToRequest(msg string) (*Request, error) {
	var (
		data        map[string]string
		requestInfo = &Request{}
	)
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		return nil, err
	}
	for k, v := range data {
		switch k {
		case "url":
			u, err := url.Parse(v)
			if err != nil {
				return nil, err
			}
			requestInfo.Url = u
		case "method":
			requestInfo.Method = v
		case "retry":
			retry, err := strconv.Atoi(v)
			if err != nil {
				retry = 3
			}
			requestInfo.Retry = retry
		case "data":
			var d map[string]interface{}
			err := json.Unmarshal([]byte(v), &d)
			if err != nil {
				return nil, err
			}
			requestInfo.Data = d
		case "json":
			var j map[string]interface{}
			err := json.Unmarshal([]byte(v), &j)
			if err != nil {
				return nil, err
			}
			requestInfo.Json = j
		case "meta":
			var j map[string]interface{}
			err := json.Unmarshal([]byte(v), &j)
			if err != nil {
				return nil, err
			}
			requestInfo.Meta = j
		case "callback":
			requestInfo.Callback = v
		}
	}
	return requestInfo, err
}
