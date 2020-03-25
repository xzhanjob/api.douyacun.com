package chat

import (
	"dyc/internal/logger"
	"dyc/internal/module/account"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pkg/errors"
	"net"
)

var (
	epoller *epoll
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func init() {
	var err error
	epoller, err = MakeEpoll()
	if err != nil {
		panic(errors.Wrap(err, "make epoll error"))
	}
	go start()
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Conn    net.Conn
	account *account.Account
}

// serveWs handles websocket requests from the peer.
func ServeWs(ctx *gin.Context) {
	if a, ok := ctx.Get("account"); ok {
		acct := a.(*account.Account)
		conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
		if err != nil {
			logger.Errorf("ws upgrade http: %v", err)
			return
		}
		client := &Client{
			Conn:    conn,
			account: acct,
		}
		if err := epoller.Add(client); err != nil {
			logger.Wrap(err, "epoll add error")
			return
		}
	}
	return
}

func (c *Client) toMap() map[string]interface{} {
	return map[string]interface{}{
		"id":   c.account.Id,
		"name": c.account.Name,
	}
}

func start() {
	for{
		clients, err := epoller.Wait()
		if err != nil {
			logger.Debugf("")
			continue
		}
		for _, client := range clients {
			message, op, err := wsutil.ReadClientData(client.Conn)
			if err != nil {
				logger.Errorf("read client data error: %v", err)
				continue
			}
			if op == ws.OpText { // 文本消息
				msg := ClientMessage{}
				if err := json.Unmarshal(message, msg); err != nil {
					logger.Errorf("json unmarshal error: %v", err)
					continue
				}
			}
		}
	}
}
