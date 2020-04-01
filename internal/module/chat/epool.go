package chat

import (
	"dyc/internal/consts"
	"dyc/internal/logger"
	"github.com/gobwas/ws/wsutil"
	"net"
	"reflect"
	"syscall"
)

type Responser interface {
	Members() []string
	Bytes() []byte
	GetChannelID() string
}

type epoll struct {
	Fd int
	// 方便通过账户id映射到文件描述符，通过文件描述符取到对应connect
	accounts    map[string]int
	connections map[int]*Client
	register    chan Client
	unregister  chan Client
	broadcast   chan Responser
	pool        *pool
}

func MakeEpoll() (*epoll, error) {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		Fd:          fd,
		accounts:    make(map[string]int),
		connections: make(map[int]*Client, 0),
		register:    make(chan Client),
		unregister:  make(chan Client),
		broadcast:   make(chan Responser),
		pool:        NewPool(128, 128),
	}, nil
}

func (e *epoll) run() {
	for {
		select {
		case client := <-e.register:
			logger.Debugf("new client register: %v", client)
			// 单点登录
			if fd, ok := e.accounts[client.account.Id]; ok {
				if other, ok := e.connections[fd]; ok {
					other.conn.Close()
					wsutil.WriteClientText(client.conn, NewTipMessage("该账号的其他链接已经关闭").Bytes())
				}
			}
			fd := websocketFd(client.conn)
			if err := syscall.EpollCtl(e.Fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Events: syscall.EPOLLIN | syscall.EPOLLOUT, Fd: int32(fd)}); err != nil {
				logger.Errorf("epoll ctl add error: %v", err)
				continue
			}
			e.accounts[client.account.Id] = fd
			e.connections[fd] = &client
		case client := <-e.unregister:
			if fd, ok := e.accounts[client.account.Id]; ok {
				if other, ok := e.connections[fd]; ok {
					if err := syscall.EpollCtl(e.Fd, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
						logger.Errorf("epoll ctl del error: %v", err)
						continue
					}
					delete(e.accounts, client.account.Id)
					delete(e.connections, fd)
					other.conn.Close()
				}
			}
		case msg := <-e.broadcast:
			logger.Debugf("new message broadcast %s", msg.Bytes())
			// 广播
			if msg.GetChannelID() == consts.GlobalChannelId {
				for id, fd := range e.accounts {
					if client, ok := e.connections[fd]; ok {
						if err := e.pool.Schedule(func() {
							client.conn.Write(msg.Bytes())
						}); err != nil {
							client.conn.Close()
							delete(e.accounts, id)
							delete(e.connections, fd)
						}
					}
				}
			} else {
				// channel聊天
				for _, id := range msg.Members() {
					if fd, ok := e.accounts[id]; ok {
						if client, ok := e.connections[fd]; ok {
							if err := e.pool.Schedule(func() {
								client.conn.Write(msg.Bytes())
							}); err != nil {
								client.conn.Close()
								delete(e.accounts, id)
								delete(e.connections, fd)
							}
						}
					}
				}
			}
		}
	}
}

func (e *epoll) Count() int {
	return len(e.connections)
}

func (e *epoll) Wait() ([]*Client, error) {
	events := make([]syscall.EpollEvent, 100)
	n, err := syscall.EpollWait(e.Fd, events, 100)
	if err != nil {
		return nil, err
	}
	connections := make([]*Client, n)
	for i := 0; i < n; i++ {
		if conn, ok  := e.connections[int(events[i].Fd)]; ok {
			connections = append(connections, conn)
		}
	}
	return connections, nil
}

func websocketFd(conn net.Conn) int {
	netConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fd := netConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
