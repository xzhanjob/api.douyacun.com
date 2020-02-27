package chat

import (
	"dyc/internal/logger"
	"encoding/json"
	"time"
)

type respType string

type response struct {
	// 类型
	RespType respType
}

const (
	RegisterRespConst respType = "REGISTER"
	SystemRespConst   respType = "SYSTEM"
	MsgRespConst      respType = "MESSAGE"
)

type RegisterResp struct {
	response
	Author *Client
}

func (r *RegisterResp) Bytes() []byte {
	data, err := json.Marshal(map[string]interface{}{
		"type":   r.RespType,
		"account": r.Author.toMap(),
	})
	if err != nil {
		logger.Errorf("msg json encode failed, %s", err)
	}
	return data
}

func NewRegisterResp(c *Client) Responser {
	return &RegisterResp{
		response: response{
			RespType: RegisterRespConst,
		},
		Author: c,
	}
}

type MsgResp struct {
	response
	Msg *Message
}

func (r *MsgResp) Bytes() []byte {
	data, err := json.Marshal(map[string]interface{}{
		"type":    r.RespType,
		"message": r.Msg.toMap(),
	})
	if err != nil {
		logger.Errorf("msg json encode failed, %s", err)
	}
	return data
}

func NewMsgResp(c *Client, msg string, t msgType) Responser {
	return &MsgResp{
		response: response{
			RespType: MsgRespConst,
		},
		Msg: &Message{
			Id:      genMsgId(),
			Content: msg,
			date:    time.Time{},
			Source:  c,
			Dest:    nil,
			MsgType: t,
		},
	}
}

func NewSystemMsg(msg string) *Message {
	return &Message{
		Id:      0,
		Content: "gst",
		date:    time.Time{},
		Source:  nil,
		Dest:    nil,
		MsgType: "",
	}
}
