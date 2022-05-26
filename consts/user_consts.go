package consts

//用户业务定义的常量
const (
	

	FOR_SEX_MALE   = "男"
	FOR_SEX_FEMALE = "女"
	FOR_SEX_BOTH   = "不限"

	CUST_SEX_MALE   = "男"
	CUST_SEX_FEMALE = "女"
	CUST_SEX_UNKNOW = "未知"

	MARRIAGE_YES = "已婚"
	MARRIAGE_NO  = "未婚"

	CUST_CARD_TYPE_ID_IDCARD   = "1" //身份证
	CUST_CARD_TYPE_ID_PASSPORT = "2" //护照
	CUST_CARD_TYPE_ID_HUIXIANG = "3" //回乡证
	CUST_CARD_TYPE_ID_TAIWAN   = "4" //台胞证
	CUST_CARD_TYPE_ID_OTHER    = "5" //其他

	WARNING_URL  = "" // 告警
	WARNING_NAME = ""
)

//枚举
var ENUM_CUST_CARD_TYPE = map[string]string{"1": "身份证", "2": "护照", "3": "回乡证", "4": "台胞证", "5": "其他"}
var ENUM_CUST_SEX = []string{"男", "女", "未知"} //客户性别
var ENUM_CUST_MARRIAGE = []string{"未婚", "已婚"}
var ENUM_FOR_SEX = []string{"男", "女", "不限"} //适应哪些性别
