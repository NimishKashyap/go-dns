package main

import (
	"encoding/binary"
	"fmt"

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
	Header DNSHeader
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
		header := newDNSHeader()
		response := header.toBytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
