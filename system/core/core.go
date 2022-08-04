package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"dsjk.com/dwggo/system/lib/helper/str"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//加载配置文件
var config_files = []string{"./config/config.yaml", "./config/secret.yaml", "./config/system.yaml"}
var config *viper.Viper

var config_cache_string = map[string]string{}
var config_cache_int = map[string]int{}
var config_cache_duration = map[string]time.Duration{}
var config_cache_float64 = map[string]float64{}

func init() {
	loadConfig()
	LogInfo("load config0:" + config_files[0])
}

func loadConfig() {
	config = viper.New()
	for index, arr := range config_files {
		config.SetConfigFile(arr)
		if index == 0 {
			config.ReadInConfig()
		} else {
			config.MergeInConfig()
		}
	}
	//port := config.GetString("app.port")
}
func GetConfigString(key string) string {
	if cache, ok := config_cache_string[key]; ok {
		return cache
	}
	str := config.GetString(key)
	if str == "" { //疑似配错，提示
		fmt.Println("GetConfigString is empty", key)
	}
	config_cache_string[key] = str
	return str
}
func GetConfigInt(key string) int {
	if cache, ok := config_cache_int[key]; ok {
		return cache
	}
	ret := config.GetInt(key)
	if ret == 0 { //疑似配错，提示
		ret_true := config.Get(key)
		if ret_true == nil {
			fmt.Println("GetConfigInt is empty", key)
		}
	}
	config_cache_int[key] = ret
	return ret
}
func GetConfigDuration(key string) time.Duration {
	if cache, ok := config_cache_duration[key]; ok {
		return cache
	}
	ret := config.GetDuration(key)
	if ret == 0 { //疑似配错，提示
		ret_true := config.Get(key)
		if ret_true == nil {
			fmt.Println("GetConfigDuration is empty", key)
		}
	}
	config_cache_duration[key] = ret
	return ret
}
func GetConfigFloat64(key string) float64 {
	if cache, ok := config_cache_float64[key]; ok {
		return cache
	}
	ret := config.GetFloat64(key)
	if ret == 0 { //疑似配错，提示
		ret_true := config.Get(key)
		if ret_true == nil {
			fmt.Println("GetConfigFloat64 is empty", key)
		}
	}
	config_cache_float64[key] = ret
	return ret
}

//判断运行环境
func IsEnvDev() bool {
	return gin.Mode() == gin.DebugMode
}
func IsEnvTest() bool {
	return gin.Mode() == gin.TestMode
}
func IsEnvPro() bool {
	return gin.Mode() == gin.ReleaseMode
}

/**
崩溃记录日志
*/
func ForPanicLog() {
	err := recover()
	if err != nil {
		LogError(fmt.Sprint(err))
		fmt.Println("panic error:", err)
		panic(err) //记录完日志再抛出去，让原生的接管
	}
}

//发送告警
func SendWarning(msg string) {
	if !IsEnvPro() || REPORT_WARNING_URL == "" {
		return
	}
	url := REPORT_WARNING_URL
	params := map[string]interface{}{}
	params["msg"] = msg
	params["from"] = REPORT_WARNING_NAME
	_, err := _PostForm(url, &params)
	if err != nil { //发送失败
		fmt.Printf("err: %v\n", err)
		loadErrorLogger().Error("发送告警异常：" + err.Error())
	}
}
func _PostForm(url_str string, params ...*map[string]interface{}) (*string, error) {
	if len(params) == 0 {
		return nil, errors.New("post param is request")
	}

	urlValues := url.Values{}
	for k, v := range *params[0] {
		urlValues.Add(k, str.ToString(v))
	}
	resp, err := http.PostForm(url_str, urlValues)
	if err != nil {
		emp_str := ""
		return &emp_str, nil
	}
	return _doResp(url_str, resp)
}
func _doResp(url_str string, resp *http.Response) (*string, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body_str := string(body)
	if resp.StatusCode != 200 {
		return nil, errors.New("resp.StatusCode is not 200,it is:" + strconv.Itoa(resp.StatusCode) + ",body:" + body_str)
	}

	return &body_str, nil
}
