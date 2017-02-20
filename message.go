package halo

import (
	"errors"
	"fmt"
)

type MessageType byte

const (
	Request  MessageType = 0x00
	Notify               = 0x01
	Response             = 0x02
	Push                 = 0x03
)

const (
	msgRouteCompressMask = 0x01
	msgTypeMask          = 0x07
	msgRouteLengthMask   = 0xFF
	msgHeadLength        = 0x03
)

var (
	routeDict = make(map[string]uint16)
	codeDict  = make(map[uint16]string)
)

var (
	ErrWrongMessageType  = errors.New("wrong message type")
	ErrInvalidMessage    = errors.New("invalid message")
	ErrRouteInfoNotFound = errors.New("route info not found in dictionary")
)

type Message struct {
	ServerType string
	ID         uint
	Data       []byte
	compressed bool
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) String() string {
	return fmt.Sprintf("Server: %s, ID: %d, Compressed: %t, BodyLength: %d",
		m.ServerType,
		m.ID,
		m.compressed,
		len(m.Data))
}

func (m *Message) Encode() ([]byte, error) {
	return Encode(m)
}

// Encode message. Different message types is corresponding to different message header,
// message types is identified by 2-4 bit of flag field. The relationship between message
// types and message header is presented as follows:
//
//   type      flag      other
//   ----      ----      -----
// request  |----000-|<message id>|<route>
// notify   |----001-|<route>
// response |----010-|<message id>
// push     |----011-|<route>
// The figure above indicates that the bit does not affect the type of message.
func Encode(m *Message) ([]byte, error) {
	return nil, nil
}

func Decode(data []byte) (*Message, error) {
	return nil, nil
}
