package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
)

type Peer struct {
	Name      string
	TcpSocket net.Conn
	PeerList  []Peer
	isLeader  bool
}

func main() {
	conn, err := net.Dial("tcp", "localhost:4339")
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	dec := gob.NewDecoder(conn)
	p := &Peer{}
	dec.Decode(p)
	fmt.Printf("Received : %+v", p)
	fmt.Println(status)
	fmt.Println(err)
}

