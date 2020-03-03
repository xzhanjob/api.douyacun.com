package chat

import (
	"dyc/internal/module/account"
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
			"id":   m.Source.account.Id,
			"name": m.Source.account.Name,
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

func WithSystemMsg(msg string) *Message {
	a := account.NewAccount()
	a.Name = "系统消息"
	a.Id = "0"
	return &Message{
		Id:      0,
		Content: msg,
		date:    time.Time{},
		Source:  &Client{account: a},
		Dest:    nil,
		MsgType: TextMsg,
	}
}

func NewDefaultMsg(c *Client, msg string) *Message {
	return &Message{
		Id:      0,
		Content: msg,
		date:    time.Time{},
		Source:  c,
		Dest:    nil,
		MsgType: TextMsg,
	}
}

func (m *Message) GetClient() *Client {
	return m.Source
}

func genMsgId() int64 {
	atomic.AddInt64(&msgId, 1)
	return msgId
}
