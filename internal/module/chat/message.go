package chat

import (
	"sync/atomic"
	"time"
)

// 消息来源

type msgType string

const (
	TextMsg msgType = "TEXT"
	ImgMsg  msgType = "IMAGE"
	FileMsg msgType = "FILE"
)

var msgId int64

type Message struct {
	// 消息id
	Id int64
	// 内容
	Content string
	// 发送时间
	date time.Time
	// 发送者
	Source *Client
	// 接受者
	Dest *Client
	// 消息类型
	MsgType msgType
}

// 实际发送结构体
func (m *Message) toMap() map[string]interface{} {
	return map[string]interface{}{
		"source": map[string]interface{}{
			"id":   m.Source.id,
			"name": m.Source.getName(),
		},
		"dest": map[string]interface{}{
		},
		"date": m.date,
		"context": map[string]interface{}{
			"id":      m.Id,
			"type":    m.MsgType,
			"content": m.Content,
		},
	}
}



func (m *Message) GetClient() *Client {
	return m.Source
}

func genMsgId() int64 {
	atomic.AddInt64(&msgId, 1)
	return msgId
}
