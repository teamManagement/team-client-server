package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HttpRequestWrapper struct {
	req    *http.Request
	Client *http.Client
}

func (h *HttpRequestWrapper) ResponseToJson(res any) error {
	b, err := h.ResponseToBytes()
	if b != nil {
		if e := json.Unmarshal(b, &res); e != nil {
			if err != nil {
				return err
			}
			return e
		}
	}
	return err
}

func (h *HttpRequestWrapper) ResponseToText() (string, error) {
	b, err := h.ResponseToBytes()
	return string(b), err
}

func (h *HttpRequestWrapper) ResponseToBytes() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	return result, h.ResponseWithHandler(func(res *http.Response) error {
		result, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("读取响应体数据失败: %s", err.Error())
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return fmt.Errorf("响应状态非正常: %d", res.StatusCode)
		}

		return nil
	})
}

func (h *HttpRequestWrapper) ResponseWithHandler(fn func(res *http.Response) error) error {
	if h.Client == nil {
		h.Client = http.DefaultClient
	}

	res, err := h.Client.Do(h.req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return fn(res)
}

type HttpRequestOption struct {
	Method      string
	Header      map[string]string
	JsonData    any
	FormData    url.Values
	contentType string
	body        io.Reader
}

func (h *HttpRequestOption) check() (*HttpRequestOption, error) {
	target := h
	if target == nil {
		target = &HttpRequestOption{}
	}

	if target.Method == "" {
		target.Method = "POST"
	}

	target.Method = strings.ToUpper(target.Method)

	if target.JsonData != nil {
		marshal, err := json.Marshal(target.JsonData)
		if err != nil {
			return nil, fmt.Errorf("序列化JSON数据失败: %s", err.Error())
		}

		target.body = bytes.NewReader(marshal)
		target.contentType = "application/json"
	} else if target.FormData != nil {
		target.body = strings.NewReader(target.FormData.Encode())
		target.contentType = "application/x-www-form-urlencoded"
	}

	return target, nil
}

func HttpRequest(url string) (request *HttpRequestWrapper, err error) {
	return HttpRequestWithOption(url, nil)
}

func HttpRequestWithOption(url string, option *HttpRequestOption) (*HttpRequestWrapper, error) {
	var (
		request *http.Request
		err     error
	)

	if option, err = option.check(); err != nil {
		return nil, err
	}

	request, err = http.NewRequest(option.Method, url, option.body)
	if err != nil {
		return nil, err
	}

	for k, v := range option.Header {
		request.Header.Set(k, v)
	}

	if option.contentType != "" {
		request.Header.Set("Content-Type", option.contentType)
	}

	return &HttpRequestWrapper{
		req: request,
	}, nil

}
