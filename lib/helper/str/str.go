package str

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TIME_YMDHIS = "2006-01-02 15:04:05" //用于格式化时间YYYY-MM-DD HH:ii:ss
const TIME_YMD = "2006-01-02"             //用于格式化时间YYYY-MM-DD

// 字符串首字母大写
func FirstToUpper(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func EqualStrInterface(str string, ele interface{}) bool {
	return ToString(ele) == str
}

//字符串转int
func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
func ToIntAny(str interface{}) int {
	return ToInt(ToString(str))
}

//字符串转int
func ToFloat32(str string) float32 {
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}
func ToFloat32Any(str interface{}) float32 {
	return ToFloat32(ToString(str))
}

//将interface{}类型的变量转成字符串
// 浮点型 3.0将会转换成字符串3, "3"
func ToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case bool:
		// newValue, _ := json.Marshal(value)
		// key = string(newValue)
		key_bool := value.(bool)
		if key_bool {
			key = "true"
		} else {
			key = "false"
		}

	case *float64:
		ft := value.(*float64)
		key = strconv.FormatFloat(*ft, 'f', -1, 64)
	case *float32:
		ft := value.(*float32)
		key = strconv.FormatFloat(float64(*ft), 'f', -1, 64)
	case *int:
		it := value.(*int)
		key = strconv.Itoa(*it)
	case *uint:
		it := value.(*uint)
		key = strconv.Itoa(int(*it))
	case *int8:
		it := value.(*int8)
		key = strconv.Itoa(int(*it))
	case *uint8:
		it := value.(*uint8)
		key = strconv.Itoa(int(*it))
	case *int16:
		it := value.(*int16)
		key = strconv.Itoa(int(*it))
	case *uint16:
		it := value.(*uint16)
		key = strconv.Itoa(int(*it))
	case *int32:
		it := value.(*int32)
		key = strconv.Itoa(int(*it))
	case *uint32:
		it := value.(*uint32)
		key = strconv.Itoa(int(*it))
	case *int64:
		it := value.(*int64)
		key = strconv.FormatInt(*it, 10)
	case *uint64:
		it := value.(*uint64)
		key = strconv.FormatUint(*it, 10)
	case *string:
		key = *(value.(*string))
	case *[]byte:
		key = string(*(value.(*[]byte)))
	case *bool:
		// newValue, _ := json.Marshal(value)
		// key = string(newValue)
		key_bool := value.(*bool)
		if *key_bool {
			key = "true"
		} else {
			key = "false"
		}
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
		//panic("interface ToString(): unsupported type: " + fmt.Sprint(value))
	}
	return key
}

func Json_decode(json_str *string) *map[string]interface{} {
	var ret_json map[string]interface{}
	json.Unmarshal([]byte(*json_str), &ret_json)
	return &ret_json
}

func Json_encode(m interface{}) string {
	newValue, _ := json.Marshal(m)
	return string(newValue)
}

func Base64_encode(str string) *string {
	str_ret := base64.StdEncoding.EncodeToString([]byte(str))
	return &str_ret
}
func Base64_decode(str string) *string {
	sDec, _ := base64.StdEncoding.DecodeString(str)
	str_ret := string(sDec)
	return &str_ret
}

//日期时间格式转换
//TimeFormat(str.TIME_YMD,"2021-12-01 12:30:30")
//return 2021-12-01
func TimeFormat(format, input_date string) string {
	t, _ := time.Parse(TIME_YMDHIS, input_date)
	return t.Format(format)
}
func GetAgeByBirthday(birthday string) int {
	t, err := time.Parse(TIME_YMD, birthday)
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - t.Year()
	if now.YearDay() >= t.YearDay() { //是否过了生日。。
		age = age + 1
	}
	if age < 0 {
		age = 0
	}
	return age
}

//匹配字符串开头，如果有特殊符号，正则可能出错
func StartWith(str *string, prefixs []string) bool {
	reg_str := strings.Join(prefixs, "|")
	m, err := regexp.MatchString(`^(`+reg_str+`)`, *str)
	if err != nil {
		panic(err.Error())
	}
	return m
}
