package myapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/lib/helper/str"
)

const (
	MYAPI_METHOD_GET       = "GET"
	MYAPI_METHOD_POST_JSON = "JSON"
	MYAPI_METHOD_POST_FORM = "POST"
	//POST请求发送文本内容，不是json。。
	// params_tmp := map[string]interface{}{}
	// params_tmp[myapi.MYAPI_POST_TEXT_CONTENT] = 文本内容
	// res, err := self.Api(url_str, myapi.MYAPI_METHOD_POST_TEXT, &params_tmp)
	MYAPI_METHOD_POST_TEXT  = "POST_TEXT"
	MYAPI_METHOD_WEBSERVICE = "WEBSERVICE"

	MYAPI_POST_TEXT_CONTENT = "_content_"
)

//API请求-回写请求和响应日志
func Api(url_str string, method string, params ...map[string]interface{}) (*string, error) {
	return ApiALLParams(url_str, method, true, true, params...)
}

//API请求-不写请求和响应日志
func ApiNoAllLog(url_str string, method string, params ...map[string]interface{}) (*string, error) {
	return ApiALLParams(url_str, method, false, false, params...)
}

//API请求-写请求日志，不写响应日志
func ApiNoResponsLog(url_str string, method string, params ...map[string]interface{}) (*string, error) {
	return ApiALLParams(url_str, method, true, false, params...)
}
func ApiALLParams(url_str string, method string, WriteRequestLog, WriteResponsLog bool, params ...map[string]interface{}) (*string, error) {
	var body *string
	var err error
	time_start := time.Now()
	defer func() {
		WriteReqResRecord(body, err, url_str, method, WriteRequestLog, WriteResponsLog, time_start, params...) //匿名函数才能获取到真正的body，否则是nil
	}()
	body, err = RESTFul(url_str, method, params...)

	return body, err
}
func WriteReqResRecord(body *string, err error, url_str string, method string, WriteRequestLog, WriteResponsLog bool, time_start time.Time, params ...map[string]interface{}) { //记录请求，返回日志
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
		time_offset := fmt.Sprintf("%.3f", (float64(time.Since(time_start)) / float64(time.Second)))
		core.LogInfoCustom("url_str:"+url_str+"	"+method+"	times:"+time_offset+"s	params:"+string(json)+"	respons:"+res+err_str, core.APILOG_FILENAME)
	}
}

func RESTFul(url_str string, method string, params ...map[string]interface{}) (*string, error) {
	switch strings.ToUpper(method) {
	case MYAPI_METHOD_GET:
		return Get(url_str, params...)
	case MYAPI_METHOD_POST_FORM:
		return PostForm(url_str, params...)
	case MYAPI_METHOD_POST_JSON:
		return PostJson(url_str, params...)
	case MYAPI_METHOD_POST_TEXT:
		return PostText(url_str, params...)
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
func Get(url_str string, params ...map[string]interface{}) (*string, error) {
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
func PostForm(url_str string, params ...map[string]interface{}) (*string, error) {
	var err error
	if len(params) == 0 {
		return nil, errors.New("post param is request")
	}
	//options := params[1] TODO

	urlValues := url.Values{}
	for k, v := range params[0] {
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
func PostJson(url_str string, params ...map[string]interface{}) (*string, error) {
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

func PostText(url_str string, params ...map[string]interface{}) (*string, error) {
	var err error
	if len(params) == 0 {
		return nil, errors.New("post param is request")
	}
	//options := params[1] TODO

	bytesData, ok := (params[0])[MYAPI_POST_TEXT_CONTENT]
	if !ok {
		return nil, errors.New("post param miss:" + MYAPI_POST_TEXT_CONTENT)
	}
	resp, err := http.Post(url_str, "application/json", bytes.NewReader([]byte(str.ToString(bytesData))))
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
		log_str := "resp.StatusCode is not 200,url:" + url_str + ",http_code:" + strconv.Itoa(resp.StatusCode) + ",body:" + body_str
		core.LogError(log_str)
		return nil, errors.New(log_str)
	}

	return &body_str, nil
}

/**
拼接get参数到url上
*/
func AppendQueryToUrl(url_str string, querys map[string]interface{}) (string, error) {
	params := url.Values{}
	url_ret, err := url.Parse(url_str)
	if err != nil {
		core.LogError("AppendQueryToUrl parse url error:" + url_str + "," + err.Error())
		return "", err
	}
	for k, v := range querys {
		params.Set(k, str.ToString(v))
	}
	if url_ret.RawQuery != "" {
		url_ret.RawQuery = url_ret.RawQuery + "&" + params.Encode()
	} else {
		url_ret.RawQuery = params.Encode()
	}

	return url_ret.String(), nil
}
