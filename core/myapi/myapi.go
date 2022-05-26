package myapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"dsjk.com/openplatform/consts"
	"dsjk.com/openplatform/core"

	"dsjk.com/openplatform/lib/helper/str"
)

const (
	MYAPI_METHOD_GET        = "GET"
	MYAPI_METHOD_POST_JSON  = "JSON"
	MYAPI_METHOD_POST_FORM  = "POST"
	MYAPI_METHOD_WEBSERVICE = "WEBSERVICE"
)

//API请求-回写请求和响应日志
func Api(url_str string, method string, params ...*map[string]interface{}) (*string, error) {
	return ApiALLParams(url_str, method, true, true, params...)
}

//API请求-不写请求和响应日志
func ApiNoAllLog(url_str string, method string, params ...*map[string]interface{}) (*string, error) {
	return ApiALLParams(url_str, method, false, false, params...)
}

//API请求-写请求日志，不写响应日志
func ApiNoResponsLog(url_str string, method string, params ...*map[string]interface{}) (*string, error) {
	return ApiALLParams(url_str, method, true, false, params...)
}
func ApiALLParams(url_str string, method string, WriteRequestLog, WriteResponsLog bool, params ...*map[string]interface{}) (*string, error) {
	var body *string
	var err error
	defer func() {
		writeReqResRecord(body, err, url_str, method, WriteRequestLog, WriteResponsLog, params...) //匿名函数才能获取到真正的body，否则是nil
	}()
	body, err = RESTFul(url_str, method, params...)

	return body, err
}
func writeReqResRecord(body *string, err error, url_str string, method string, WriteRequestLog, WriteResponsLog bool, params ...*map[string]interface{}) { //记录请求，返回日志
	if WriteRequestLog {
		json, _ := json.Marshal(params)
		var res string
		err_str := ""
		if WriteResponsLog == false {
			res = "*body not recorded"
		} else if body == nil {
			res = "*body is nil"
		} else {
			res = *body
		}
		if err != nil {
			err_str = ",error:" + err.Error()
		} else {
			err_str = ",error:"
		}
		core.LogInfoCustom("url_str:"+url_str+"	"+method+"	params:"+string(json)+"	respons:"+res+err_str, consts.LOGFILE_MYAPI)
	}
}

func RESTFul(url_str string, method string, params ...*map[string]interface{}) (*string, error) {
	switch strings.ToUpper(method) {
	case MYAPI_METHOD_GET:
		return Get(url_str, params...)
	case MYAPI_METHOD_POST_FORM:
		return PostForm(url_str, params...)
	case MYAPI_METHOD_POST_JSON:
		return PostJson(url_str, params...)
	}

	return nil, errors.New("unsupported type:" + method)
}

/**
GET
parmas 参数
parmas[0] 请求参数
params[1] options
options {
	"timeout/content_type/ret_is_json" //TODO
}
*/
func Get(url_str string, params ...*map[string]interface{}) (*string, error) {
	var err error
	if len(params) > 0 {
		url_str, err = AppendQueryToUrl(url_str, params[0])
	}
	//options := params[1] TODO

	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url_str)
	return doResp(err, url_str, resp)
}

/**
POST -form
parmas 参数
parmas[0] 请求参数
params[1] options
options {
	"timeout/content_type/ret_is_json" //TODO
}
*/
func PostForm(url_str string, params ...*map[string]interface{}) (*string, error) {
	var err error
	if len(params) == 0 {
		return nil, errors.New("post param is request")
	}
	//options := params[1] TODO

	urlValues := url.Values{}
	for k, v := range *params[0] {
		urlValues.Add(k, str.ToString(v))
	}
	resp, err := http.PostForm(url_str, urlValues)
	return doResp(err, url_str, resp)
}

/**
POST-json
parmas 参数
parmas[0] 请求参数
params[1] options
options {
	"timeout/content_type/ret_is_json" //TODO
}
*/
func PostJson(url_str string, params ...*map[string]interface{}) (*string, error) {
	var err error
	if len(params) == 0 {
		return nil, errors.New("post param is request")
	}
	//options := params[1] TODO

	bytesData, err := json.Marshal(params[0])
	if err != nil {
		core.LogError(url_str + "," + err.Error())
		return nil, err
	}
	resp, err := http.Post(url_str, "application/json", bytes.NewReader(bytesData))
	return doResp(err, url_str, resp)
}

func doResp(err error, url_str string, resp *http.Response) (*string, error) {
	if err != nil {
		core.LogError(url_str + "," + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.LogError(err.Error())
		return nil, err
	}
	body_str := string(body)
	if resp.StatusCode != 200 {
		core.LogError("resp.StatusCode is not 200," + url_str + "," + strconv.Itoa(resp.StatusCode))
		return nil, errors.New("resp.StatusCode is not 200,it is:" + strconv.Itoa(resp.StatusCode) + ",body:" + body_str)
	}

	return &body_str, nil
}

/**
拼接get参数到url上
*/
func AppendQueryToUrl(url_str string, querys *map[string]interface{}) (string, error) {
	params := url.Values{}
	url_ret, err := url.Parse(url_str)
	if err != nil {
		core.LogError("AppendQueryToUrl parse url error:" + url_str + "," + err.Error())
		return "", err
	}
	for k, v := range *querys {
		params.Set(k, str.ToString(v))
	}
	if url_ret.RawQuery != "" {
		url_ret.RawQuery = url_ret.RawQuery + "&" + params.Encode()
	} else {
		url_ret.RawQuery = params.Encode()
	}

	return url_ret.String(), nil
}
