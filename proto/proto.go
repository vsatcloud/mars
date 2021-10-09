package proto

//Pagination 分页
type Pagination struct {
	//页码 从1开始
	Page int `form:"page" default:"1" json:"page"`
	//兼容字段-勿用，请使用size字段
	Limit int `form:"limit" default:"10" json:"limit"`
	//每页的数量
	Size int `form:"size" default:"10" json:"size"`
}

func (p *Pagination) Offset() int {
	if p.Page <= 0 {
		return -1 //cancel offset
	}

	if p.Size != 0 && p.Size != 10 {
		p.Limit = p.Size
	}

	return (p.Page - 1) * p.Limit
}

func (p *Pagination) Limiter() int {
	if p.Limit <= 0 {
		return -1
	}

	return p.Limit
}

// Result represents HTTP response body.
type Result struct {
	Code    int         `json:"code"`    // return code, 0 for succ
	Message string      `json:"message"` // message
	Data    interface{} `json:"data"`    // data object
	Detail  string      `json:"detail"`
}

type ResultList struct {
	Code    int    `json:"code"`    // return code, 0 for succ
	Message string `json:"message"` // message
	Data    struct {
		Count int64         `json:"count"`
		Items []interface{} `json:"items"`
		Pagination
	} `json:"data"` // data object
	Detail string `json:"detail"`
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
	CodeSuccess:          "Success",
	CodeErrSystem:        "系统内部错误，请联系相关人员",
	CodeErrParams:        "参数错误",
	CodeFailedAuthVerify: "身份验证失败，请重新登录",
	CodeNoPerm:           "没有访问权限，请联系相关人员",
	CodeNotFound:         "资源不存在",
	CodeTokenExpired:     "身份已过期，请重新登录",
}
