package model

//Ademo控制器Query方法的请求参数
type ArgAdemoQueryInModel struct {
	Name string `form:"name" json:"name"`
}

//Ademo控制器Query方法的返回参数
type ArgAdemoQueryOutModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
