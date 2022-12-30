package remoteserver

import (
	"encoding/json"
	"fmt"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"strings"
	"team-client-server/tools"
)

func RequestWebService(url string) error {
	return RequestWebServiceWithData(url, nil, nil)
}

func RequestWebServiceToRawReq(url string, data any) (*tools.HttpRequestWrapper, error) {
	option := &tools.HttpRequestOption{
		Method:   "POST",
		JsonData: data,
		Header: map[string]string{
			"_t":         Token(),
			"_a":         LoginIp(),
			"User-Agent": "teamManageLocal",
		},
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	url = LocalWebServerAddress + url

	return tools.HttpRequestWithOption(url, option)
}

func RequestWebServiceWithData(url string, data any, res any) error {
	option := &tools.HttpRequestOption{
		Method:   "POST",
		JsonData: data,
		Header: map[string]string{
			"_t":         Token(),
			"_a":         LoginIp(),
			"User-Agent": "teamManageLocal",
		},
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	url = LocalWebServerAddress + url

	req, err := tools.HttpRequestWithOption(url, option)
	if err != nil {
		return err
	}

	var httpResult *ginmiddleware.HttpResult
	if err = req.ResponseToJson(&httpResult); err != nil && httpResult == nil {
		return err
	} else if err != nil && httpResult != nil {
		if httpResult.Result == nil {
			return err
		}

	}

	if httpResult.Error {
		return fmt.Errorf("%s: %s", httpResult.Code, httpResult.Msg)
	}

	if res == nil {
		return err
	}

	marshal, err := json.Marshal(httpResult.Result)
	if err != nil {
		return fmt.Errorf("解析响应体内容失败: %s", err.Error())
	}

	return json.Unmarshal(marshal, &res)
}

func RequestWebServiceWithResponse(url string, res any) error {
	return RequestWebServiceWithData(url, nil, res)
}
