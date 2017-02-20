package halo

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type networkStatus byte

const (
	_ networkStatus = iota
	StatusStart
	StatusHandshake
	StatusWorking
	StatusClosed
)

var (
	ErrRPCLocal     = errors.New("RPC object must location in different server type")
	ErrSidNotExists = errors.New("sid not exists")
)

// Agent corresponding a user, used for store raw socket information
// only used in package internal, can not accessible by other package
type Agent struct {
	id       int64
	socket   net.Conn
	session  *Session
	lastTime int64 // last heartbeat unix time stamp
}

// Create new Agent instance
func newAgent(conn net.Conn) *Agent {
	a := &Agent{
		socket:   conn,
		lastTime: time.Now().Unix()}
	s := NewSession(a)
	a.session = s
	a.id = s.ID
	return a
}

// String, implementation for Stringer interface
func (a *Agent) String() string {
	return fmt.Sprintf("id: %d, remote address: %s, last time: %d",
		a.id,
		a.socket.RemoteAddr().String(),
		a.lastTime)
}

func (a *Agent) Heartbeat() {
	a.lastTime = time.Now().Unix()
}

func (a *Agent) Close() {
	DefaultNetService.CloseSession(a.session)
	a.socket.Close()
}

func (a *Agent) ID() int64 {
	return a.id
}

func (a *Agent) Session() *Session {
	return a.session
}

func (a *Agent) Socket() net.Conn {
	return a.socket
}

func (a *Agent) Send(data []byte) error {
	_, err := a.socket.Write(data)
	return err
}

func (a *Agent) Push(session *Session, route string, v interface{}) error {
	data, err := serializeOrRaw(v)
	if err != nil {
		return err
	}
	return DefaultNetService.Push(session, route, data)
}

// Response message to session
func (a *Agent) Response(session *Session, v interface{}) error {
	data, err := serializeOrRaw(v)
	if err != nil {
		return err
	}

	return DefaultNetService.Response(session, data)
}
