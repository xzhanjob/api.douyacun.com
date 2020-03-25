package chat

import (
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
					other.Conn.Close()
				}
			}
			fd := websocketFd(client.Conn)
			e.accounts[client.account.Id] = fd
			e.connections[fd] = &client
		case client := <-e.unregister:
			// 单点登录
			if fd, ok := e.accounts[client.account.Id]; ok {
				if other, ok := e.connections[fd]; ok {
					delete(e.accounts, client.account.Id)
					delete(e.connections, fd)
					other.Conn.Close()
				}
			}
		}
	}
}

func (e *epoll) Add(c *Client) error {
	fd := websocketFd(c.Conn)
	if err := syscall.EpollCtl(e.Fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Events: syscall.EPOLLIN | syscall.EPOLLOUT, Fd: int32(fd)}); err != nil {
		return err
	}
	e.connections[fd] = c
	return nil
}

func (e *epoll) Remove(c *Client) error {
	fd := websocketFd(c.Conn)
	if err := syscall.EpollCtl(e.Fd, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
		return err
	}
	e.Lock.Lock()
	defer e.Lock.Unlock()
	delete(e.connections, fd)
	return nil
}

func (e *epoll) Count(c Client) int {
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
