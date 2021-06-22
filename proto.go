package mars

//Pagination 分页
type Pagination struct {
	Page  int `form:"page" default:"1" json:"page"`
	Limit int `form:"limit" default:"10" json:"limit"`
}

// Result represents HTTP response body.
type Result struct {
	Code    int         `json:"code"`    // return code, 0 for succ
	Message string      `json:"message"` // message
	Data    interface{} `json:"data"`    // data object
	Detail  string      `json:"detail"`
}

const (
	CodeSuccess          = 0
	CodeErrSystem        = 10001 //系统错误
	CodeErrParams        = 10002 //参数错误
	CodeErrLogic         = 10003 //逻辑错误
	CodeFailedAuthVerify = 10004 //身份验证失败
	CodeRecordExists     = 10005 //记录已存在
	CodeNoPerm           = 10006 //没有权限
	CodeNotFound         = 10007 //没有找到
	CodeTokenExpired     = 10008 //身份过期
)

var CodeMsg = map[int]string{
	CodeSuccess:      "Success",
	CodeErrSystem:    "系统错误",
	CodeErrParams:    "参数错误",
	CodeNoPerm:       "没有权限",
	CodeNotFound:     "资源不存在",
	CodeTokenExpired: "身份已过期，请重新登录",
}
