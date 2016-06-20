package gostun

import (
	"encoding/binary"
	"errors"
	"net"
)

//StunMagic is the constant we expect from all STUN responses/requests in the Magic field
const StunMagic = 0x2112A442

//BindingX is the byte indicator of the STUN message type
const (
	BindingRequest         = 0x0001
	BindingSuccessResponse = 0x0101
)

//AttributeX is the byte indicator of the STUN attribute type
const (
	AttributeReserved         = 0x0000
	AttributeMappedAddress    = 0x0001
	AttributeResponseAddress  = 0x0002
	AttributeChangeAddress    = 0x0003
	AttributeSourceAddress    = 0x0004
	AttributeChangedAddress   = 0x0005
	AttributeUsername         = 0x0006
	AttributePassword         = 0x0007
	AttributeMessageIntegrity = 0x0008
	AttributeErrorCode        = 0x0009
	AttributeUnknown          = 0x000A
	AttributeReflectedFrom    = 0x000B
	AttributeRealm            = 0x0014
	AttributeNonce            = 0x0015
	AttributeXORMappedAddress = 0x0020
)

//ErrX are the errors to be expected during message handling
var (
	ErrInvalidRequest = errors.New("Invalid STUN request")
)

//Message holds the information about a STUN Message
type Message struct {
	MessageType   uint16
	MessageLength uint16
	Magic         uint32
	TID           []byte
	Attributes    map[uint16][]byte
}

//UnMarshal creates a Message object from data received by the STUN server
func UnMarshal(data []byte) (*Message, error) {
	length := len(data)
	if length < 20 {
		return nil, ErrInvalidRequest
	}

	msg := new(Message)

	//parse the header
	msg.MessageType = binary.BigEndian.Uint16(data[0:2])
	//check to make sure this is a binding request
	if msg.MessageType != BindingRequest {
		return nil, ErrInvalidRequest
	}
	msg.MessageLength = binary.BigEndian.Uint16(data[2:4])
	msg.Magic = binary.BigEndian.Uint32(data[4:8])
	//error on invalid Magic number
	if msg.Magic != StunMagic {
		return nil, ErrInvalidRequest
	}

	msg.TID = data[8:20]

	//if we have leftover data, parse as attributes
	if length > 20 {
		msg.Attributes = make(map[uint16][]byte)
		i := 20
		for i < length {
			attrType := binary.BigEndian.Uint16(data[i : i+2])
			attrLength := binary.BigEndian.Uint16(data[i+2 : i+4])
			i += 4 + int(attrLength)
			msg.Attributes[attrType] = data[i-int(attrLength) : i]
			if pad := int(attrLength) % 4; pad > 0 {
				i += 4 - pad
			}
		}
		//recover here to catch any index errors
		if recover() != nil {
			return nil, ErrInvalidRequest
		}
	}

	return msg, nil
}

//Marshal transforms a message into a byte array
func Marshal(m *Message) ([]byte, error) {
	result := make([]byte, 576)
	//first do the header
	binary.BigEndian.PutUint16(result[:2], m.MessageType)
	binary.BigEndian.PutUint32(result[4:8], m.Magic)
	result = append(result[:8], m.TID...)

	//now we do the attributes
	if m.Attributes != nil {
		i := 20
		for t, v := range m.Attributes {
			length := len(v)
			binary.BigEndian.PutUint16(result[i:i+2], t)
			binary.BigEndian.PutUint16(result[i+2:i+4], uint16(len(v)))
			result = append(result[:i+4], v...)
			i += 4 + length
			//if we need to pad, do so
			if pad := length % 4; pad > 0 {
				result = append(result, make([]byte, 4-pad)...)
				i += 4 - pad
			}
		}
		binary.BigEndian.PutUint16(result[2:4], uint16(i-20))
	}
	return result, nil
}

func addMappedAddress(m *Message, raddr *net.UDPAddr) {
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(raddr.Port))
	addr := raddr.IP.To4()
	m.Attributes[AttributeMappedAddress] = append([]byte{0, 1}, append(port, addr...)...)
}

func addXORMappedAddress(m *Message, raddr *net.UDPAddr) {

	addr := raddr.IP.To4()
	port := uint16(raddr.Port)
	xbytes := xorAddress(port, addr)
	m.Attributes[AttributeXORMappedAddress] = append([]byte{0, 1}, xbytes...)

}

func xorAddress(port uint16, addr []byte) []byte {

	xport := make([]byte, 2)
	xaddr := make([]byte, 4)
	binary.BigEndian.PutUint16(xport, port^uint16(StunMagic>>16))
	binary.BigEndian.PutUint32(xaddr, binary.BigEndian.Uint32(addr)^StunMagic)
	return append(xport, xaddr...)

}
