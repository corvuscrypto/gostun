package gostun

import (
	"encoding/binary"
	"errors"
)

//Stun is the constant we expect from all STUN responses/requests in the Magic field
const STUN_MAGIC = 0x2112A442

var (
	ErrInvalidRequest = errors.New("Invalid STUN request")
)

//Message holds the information about a STUN Message
type Message struct {
	MessageType   uint16
	MessageLength uint16
	Magic         uint32
	TID           []byte
}

//NewMessage creates a Message object from data received by the STUN server
func NewMessage(data []byte) (*Message, error) {
	msg := new(Message)

	//parse the header
	msg.MessageType = binary.BigEndian.Uint16(data[0:2])
	msg.MessageLength = binary.BigEndian.Uint16(data[2:4])
	msg.Magic = binary.BigEndian.Uint32(data[4:8])
	if msg.Magic != STUN_MAGIC {
		return nil, ErrInvalidRequest
	}
	msg.TID = data[8:20]
	return msg, nil
}
