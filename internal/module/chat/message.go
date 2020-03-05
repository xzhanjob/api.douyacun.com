package chat

import (
	"bytes"
	"dyc/internal/helper"
	"dyc/internal/module/account"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// 消息来源

type msgType string

const (
	TextMsg   msgType = "TEXT"
	ImgMsg    msgType = "IMAGE"
	FileMsg   msgType = "FILE"
	SystemMsg msgType = "SYSTEM"
	TipMsg    msgType = "TIP"
)

type ServerMessage struct {
	// 消息id
	Id string `json:"id"`
	// 时间
	Date time.Time `json:"date"`
	// 发送者
	Sender *account.Account `json:"sender"`
	// 类型
	Type msgType `json:"type"`
	// 内容
	Content string `json:"content"`
	// channel
	ChannelId string `json:"channel_id"`
}

type ClientMessage struct {
	ChannelId string `json:"channelId"`
	Content   string `json:"content"`
}

func NewSystemMsg(msg string, channelId string) *ServerMessage {
	a := account.NewAccount()
	a.Name = "系统消息"
	a.Id = "0"
	m := &ServerMessage{
		Content:   msg,
		Date:      time.Time{},
		Sender:    a,
		Type:      SystemMsg,
		ChannelId: channelId,
	}
	m.Id = m.GenId()
	return m
}

func NewDefaultMsg(c *Client, msg string, channelId string) *ServerMessage {
	m := &ServerMessage{
		Content:   msg,
		Sender:    c.account,
		Type:      TextMsg,
		Date:      time.Now(),
		ChannelId: channelId,
	}
	m.Id = m.GenId()
	return m
}

func NewTipMessage(msg string) *ServerMessage {
	m := &ServerMessage{
		Content: msg,
		Date:    time.Time{},
		Sender:  account.NewSystemAccount(),
		Type:    TipMsg,
	}
	m.Id = m.GenId()
	return m
}

func (m *ServerMessage) GenId() string {
	var buf bytes.Buffer
	buf.WriteString(m.Date.String())
	buf.WriteString(m.Sender.Id)
	buf.WriteString(m.Content)
	return helper.Md532(buf.Bytes())
}

func (m *ServerMessage) Bytes() []byte {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		panic(errors.Wrap(err, "message json encode error"))
	}
	return buf.Bytes()
}

func (m *ServerMessage) Members() []string {
	ids, err := ChannelMembers.MembersIds(m.ChannelId)
	if err != nil {
		return nil
	}
	return ids
}

func (m *ServerMessage) GetChannelID() string {
	return m.ChannelId
}
