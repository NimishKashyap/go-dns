package main

import (
	"encoding/binary"
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
)

// header stucture
/*
   	                              1  1  1  1  1  1
   	0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                      ID                       |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    QDCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    ANCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    NSCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    ARCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DNSHeader struct {
	ID uint16
	// Flags contain all from opcode to everything that is in min
	FLAGS   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type Message struct {
	Header   DNSHeader
	Question Question
	Answer   Answer
}

func encodeDomains(domains []string) []byte {
	encoding := []byte{}

	for _, domain := range domains {
		labels := strings.Split(domain, ".")
		for _, label := range labels {
			encoding = append(encoding, byte(len(label)))
			encoding = append(encoding, []byte(label)...)
		}
	}
	encoding = append(encoding, '\x00')
	fmt.Println(encoding)

	return encoding
}

func newDNSHeader() *DNSHeader {
	return &DNSHeader{
		ID: 1234,
		// set QR to 1,
		FLAGS:   0x8000,
		QDCOUNT: 0,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
}

func (h *DNSHeader) toBytes() []byte {
	buf := make([]byte, 12)

	binary.BigEndian.PutUint16(buf[0:2], h.ID)
	binary.BigEndian.PutUint16(buf[2:4], h.FLAGS)
	binary.BigEndian.PutUint16(buf[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCOUNT)

	return buf
}

type (
	CLASS uint16
	TYPE  uint16
)

const (
	TYPE_A     TYPE = 1
	TYPE_NS    TYPE = 2
	TYPE_MD    TYPE = 3
	TYPE_MF    TYPE = 4
	TYPE_CNAME TYPE = 5
	TYPE_SOA   TYPE = 6
	TYPE_MB    TYPE = 7
	TYPE_MG    TYPE = 8
	TYPE_MR    TYPE = 9
	TYPE_NULL  TYPE = 10
	TYPE_WKS   TYPE = 11
	TYPE_PTR   TYPE = 12
	TYPE_HINFO TYPE = 13
	TYPE_MINFO TYPE = 14
	TYPE_MX    TYPE = 15
	TYPE_TXT   TYPE = 16
)

const (
	CLASS_IN CLASS = 1
	CLASS_CS CLASS = 2
	CLASS_CH CLASS = 3
	CLASS_HS CLASS = 4
)

type Question struct {
	QNAME  []byte
	QTYPE  TYPE
	QCLASS CLASS
}

func NewQuestion() *Question {
	return &Question{
		QNAME:  []byte{},
		QTYPE:  TYPE_A,
		QCLASS: CLASS_IN,
	}
}

func (q *Question) toBytes() []byte {
	buf := make([]byte, 4+len(q.QNAME))
	copy(buf[0:], q.QNAME)
	binary.BigEndian.PutUint16(buf[len(q.QNAME):len(q.QNAME)+2], uint16(q.QTYPE))
	binary.BigEndian.PutUint16(buf[len(q.QNAME)+2:len(q.QNAME)+4], uint16(q.QCLASS))

	return buf
}

func (m *Message) encodeDomains(domains []string) {
	for _, domain := range domains {
		labels := strings.Split(domain, ".")
		for _, label := range labels {
			m.Question.QNAME = append(m.Question.QNAME, byte(len(label)))
			m.Question.QNAME = append(m.Question.QNAME, label...)
		}
	}
	m.Question.QNAME = append(m.Question.QNAME, '\x00')
	m.Header.QDCOUNT = uint16(len(domains))

	fmt.Println(m.Question.QNAME)
}

func (m *Message) toBytes() []byte {
	buf := make([]byte, 0)
	buf = append(buf, m.Header.toBytes()...)
	buf = append(buf, m.Question.toBytes()...)
	buf = append(buf, m.Answer.toBytes()...)

	return buf
}

func buildSampleResponse() []byte {
	header := newDNSHeader()
	question := NewQuestion()
	answer := NewAnswer()
	message := Message{
		Header:   *header,
		Question: *question,
		Answer:   *answer,
	}
	message.Question.QNAME = encodeDomains([]string{"codecrafters.io"})

	message.Answer.Name = encodeDomains([]string{"codecrafters.io"})
	message.Answer.RDLENGTH = 4
	message.Answer.RDATA = []byte{127, 0, 0, 1}

	message.Header.QDCOUNT = uint16(1)
	message.Header.ANCOUNT = uint16(1)

	return message.toBytes()
}

type Answer struct {
	Name     []byte
	Type     TYPE
	Class    CLASS
	TTL      int32
	RDLENGTH uint16
	RDATA    []byte
}

func NewAnswer() *Answer {
	return &Answer{
		Name:     []byte{},
		Type:     TYPE_A,
		Class:    CLASS_IN,
		TTL:      60,
		RDLENGTH: 0,
		RDATA:    []byte{},
	}
}

func (a *Answer) toBytes() []byte {
	buf := make([]byte, 10+len(a.Name)+len(a.RDATA))

	copy(buf[0:], a.Name)
	binary.BigEndian.PutUint16(buf[len(a.Name):len(a.Name)+2], uint16(a.Type))
	binary.BigEndian.PutUint16(buf[len(a.Name)+2:len(a.Name)+4], uint16(a.Class))
	binary.BigEndian.PutUint16(buf[len(a.Name)+4:len(a.Name)+8], uint16(a.TTL))
	binary.BigEndian.PutUint16(buf[len(a.Name)+8:len(a.Name)+10], uint16(a.RDLENGTH))

	copy(buf[len(a.Name)+10:], a.RDATA)

	return buf
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		// response := []byte{}
		// header := newDNSHeader()
		// response := header.toBytes()
		//
		response := buildSampleResponse()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
