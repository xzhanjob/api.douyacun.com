package chat

import (
	"sync"
	"syscall"
)

type epoll struct {
	Fd        int
	Conntions []Client
	Lock      *sync.RWMutex
}

func MakeEpoll() (*epoll, error) {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		Fd:        fd,
		Conntions: make([]Client, 0),
		Lock:      nil,
	}, nil
}

func (e *epoll) Add(c Client) error {
	
	syscall.EpollCtl(e.Fd, )
}

func websocketFd(conn )  {
	
}
