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
	hub     *epoll
)

func init() {
	var err error
	hub, err = MakeEpoll()
	if err != nil {
		panic(errors.Wrap(err, "make epoll error"))
	}
	go hub.run()
	go start()
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	conn    net.Conn
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
		client := Client{
			conn:    conn,
			account: acct,
		}
		hub.register <- client
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
	for {
		clients, err := hub.Wait()
		if err != nil {
			logger.Debugf("epoll wait %v", err)
			continue
		}
		for _, client := range clients {
			message, op, err := wsutil.ReadClientData(client.conn)
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
				hub.broadcast <- NewDefaultMsg(client, msg.Content, msg.ChannelId)
			}
		}
	}
}
