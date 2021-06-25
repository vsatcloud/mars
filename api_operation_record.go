package mars

import (
	"github.com/rs/zerolog/log"
	"github.com/vsatcloud/mars/models"
	"github.com/vsatcloud/mars/proto"
)

type ApiOperationRecordListParams struct {
	proto.Pagination
}

type OperationRecordListItem struct {
	ID           uint   `json:"id"`
	Ip           string `json:"ip"`            //请求ip
	Method       string `json:"method"`        //请求方法
	Path         string `json:"path"`          //请求路径
	Status       int    `json:"status"`        //请求状态
	Latency      string `json:"latency"`       //延迟
	Agent        string `json:"agent"`         //代理
	ErrorMessage string `json:"error_message"` //错误信息
	Body         string `json:"body"`          //请求Body
	Resp         string `json:"resp"`          //响应Body
	UserID       uint   `json:"user_id"`       //用户id
	CreatedAt    int64  `json:"created_at"`
}

type OperationRecordListData struct {
	Items []OperationRecordListItem `json:"items"`
	Count int64                     `json:"count"`
}

func ApiOperationRecordList(c *Context) {
	defer c.ResponseJson()

	var params ApiOperationRecordListParams
	err := c.BindParams(&params)
	if err != nil {
		log.Warn().Err(err).Msg("")
		c.SystemError(err)
		return
	}

	var args models.OperationRecordListArgs
	args.Pagination = params.Pagination
	list, count, err := models.OperationRecordList(&args)
	if err != nil {
		c.SystemError(err)
		return
	}

	var data OperationRecordListData
	data.Count = count
	for _, cell := range list {
		var item OperationRecordListItem
		item.ID = cell.ID
		item.Ip = cell.Ip
		item.Method = cell.Method
		item.Path = cell.Path
		item.Status = cell.Status
		item.Latency = cell.Latency.String()
		item.Agent = cell.Agent
		item.ErrorMessage = cell.ErrorMessage
		item.Body = cell.Body
		item.Resp = cell.Resp
		item.UserID = cell.UserID
		item.CreatedAt = cell.CreatedAt.Unix()

		data.Items = append(data.Items, item)
	}

	c.SetData(data)

}
