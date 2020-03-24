package chat

import (
	"net"
	"reflect"
	"sync"
	"syscall"
)

type epoll struct {
	Fd          int
	Connections map[int]*Client
	Lock        *sync.RWMutex
}

func MakeEpoll() (*epoll, error) {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		Fd:          fd,
		Connections: make(map[int]*Client, 0),
		Lock:        nil,
	}, nil
}

func (e *epoll) Add(c *Client) error {
	fd := websocketFd(c.Conn)
	if err := syscall.EpollCtl(e.Fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Events: syscall.EPOLLIN | syscall.EPOLLOUT, Fd: int32(fd)}); err != nil {
		return err
	}
	e.Lock.Lock()
	defer e.Lock.Unlock()
	e.Connections[fd] = c
	return nil
}

func (e *epoll) Remove(c *Client) error {
	fd := websocketFd(c.Conn)
	if err := syscall.EpollCtl(e.Fd, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
		return err
	}
	e.Lock.Lock()
	defer e.Lock.Unlock()
	delete(e.Connections, fd)
	return nil
}

func (e *epoll) Count(c Client) int {
	return len(e.Connections)
}

func (e *epoll) Wait() ([]*Client, error) {
	events := make([]syscall.EpollEvent, 100)
	n, err := syscall.EpollWait(e.Fd, events, 100)
	if err != nil {
		return nil, err
	}
	connections := make([]*Client, n)
	for i := 0; i < n; i++ {
		conn := e.Connections[int(events[i].Fd)]
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
