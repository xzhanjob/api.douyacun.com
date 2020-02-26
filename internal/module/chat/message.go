package chat

import (
	"dyc/internal/logger"
	"encoding/json"
	"sync/atomic"
	"time"
)

// 消息来源
type source string

const (
	RegisterSource source = "REGISTER"
	SystemSource   source = "SYSTEM"
	ChatSource     source = "CHAT"
)

var msgId int64

type Message struct {
	// 客户端id
	Id     int64
	Client *Client
	Msg    string
	Source source
	date   time.Time
}

func NewRegisterResp(c *Client, msg string) *Message {
	return NewMessage(c, RegisterSource, msg)
}

func NewSystemMsg(c *Client, msg string) *Message {
	return NewMessage(c, SystemSource, msg)
}

func NewChatMsg(c *Client, msg string) *Message {
	return NewMessage(c, ChatSource, msg)
}

func NewMessage(c *Client, t source, msg string) *Message {
	return &Message{
		Id:     genMsgId(),
		Client: c,
		Msg:    msg,
		Source: t,
		date:   time.Now(),
	}
}

func (m *Message) GetMsg() []byte {
	r := map[string]interface{}{
		"id":      m.Id,
		"author":  m.Client.getName(),
		"source":  m.Source,
		"date":    m.date,
		"message": m.Msg,
	}
	data, err := json.Marshal(r)
	if err != nil {
		logger.Errorf("msg json encode failed, %s", err)
	}
	return data
}

func (m *Message) GetClient() *Client {
	return m.Client
}

func genMsgId() int64 {
	atomic.AddInt64(&msgId, 1)
	return msgId
}
