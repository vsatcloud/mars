package models

import (
	"time"

	"github.com/vsatcloud/mars/proto"

	"gorm.io/gorm"
)

// 如果含有time.Time 请自行import time包
type OperationRecord struct {
	gorm.Model
	Ip           string        //请求ip
	Method       string        //请求方法
	Path         string        //请求路径
	Status       int           //请求状态
	Latency      time.Duration //延迟
	Agent        string        //代理
	ErrorMessage string        //错误信息
	Body         string        //请求Body
	Resp         string        //响应Body
	UserID       uint          //用户id
}

func (o *OperationRecord) Save() error {
	return db.Save(o).Error
}

type OperationRecordListArgs struct {
	proto.Pagination
}

func OperationRecordList(args *OperationRecordListArgs) (list []OperationRecord, count int64, err error) {
	query := db.Model(&OperationRecord{})
	query = query.Count(&count).Offset(Offset(args.Page, args.Limit)).Limit(Limit(args.Limit)).Scan(&list)
	err = query.Error
	return
}
