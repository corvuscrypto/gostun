package gostun

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

//Server is the main struct that contains the connection and registry information
type Server struct {
	Port       int
	Registry   *Registry
	connection *net.UDPConn
}

func (s *Server) serve() {
	for {
		var buf = make([]byte, 1024)
		size, remoteAddr, err := s.connection.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size])
	}
}

//Serve initiates a UDP connection that listens on any port for incoming data
func (s *Server) Serve() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.connection, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	go s.serve()
}

//NewServer conveniently creates a new server from the given port
func NewServer(port int) *Server {
	ret := new(Server)
	ret.Port = port
	ret.Registry = new(Registry)
	ret.Registry.mappings = make(map[string]*Client)
	return ret
}

func (s *Server) handleData(raddr *net.UDPAddr, data []byte) {
	msg, err := UnMarshal(data)
	if err != nil {
		return
	}
	respMsg := new(Message)
	respMsg.MessageType = BindingSuccessResponse
	respMsg.Magic = StunMagic
	respMsg.TID = msg.TID
	respMsg.Attributes = make(map[uint16][]byte)

	addXORMappedAddress(respMsg, raddr)
	// addMappedAddress(respMsg, raddr)

	response, err := Marshal(respMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	//send response
	_, err = s.connection.WriteToUDP(response, raddr)
	if err != nil {
		fmt.Println(err)
	}
}
