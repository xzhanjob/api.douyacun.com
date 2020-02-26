package chat

import "sync/atomic"

type msgType string

const (
	RegisterMsg msgType = "REGISTER"
	SystemMsg   msgType = "SYSTEM"
	ChatMsg     msgType = "CHAT"
)

var msgId int64

type Message struct {
	// 客户端id
	Client *Client
	Msg    []byte
	MsgType msgType
}

func (m *Message) GetMsg() []byte {
	return m.Msg
}

func (m *Message) GetClient() *Client {
	return m.Client
}

func (m *Message) GetMsgId() int64 {
	atomic.AddInt64(&msgId, 1)
	return msgId
}
