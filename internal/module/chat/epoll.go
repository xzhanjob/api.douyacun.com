package chat

import (
	"dyc/internal/consts"
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
	}, nil
}

func (e *epoll) run() {
	for {
		select {
		case client := <- e.register:
			// 单点登录
			if fd, ok := e.accounts[client.account.Id]; ok {
				if other, ok := e.connections[fd]; ok {
					other.conn.Close()
				}
			}
			fd := websocketFd(client.conn)
			e.accounts[client.account.Id] = fd
			e.connections[fd] = &client
		case client := <-e.unregister:
			if fd, ok := e.accounts[client.account.Id]; ok {
				if other, ok := e.connections[fd]; ok {
					delete(e.accounts, client.account.Id)
					delete(e.connections, fd)
					other.conn.Close()
				}
			}
		case msg := <-e.broadcast:
			// 广播
			if msg.GetChannelID() == consts.GlobalChannelId {
				for id, fd := range e.accounts {
					if client, ok := e.connections[fd]; ok {
						select {
						case client.send <- msg.Bytes():
						default:
							close(client.send)
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
							select {
							case client.send <- msg.Bytes():
							default:
								close(client.send)
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
		conn := e.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func websocketFd(conn net.Conn) int {
	netConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fd := netConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
