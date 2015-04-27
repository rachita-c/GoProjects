package main

import (
        "encoding/gob"
        "fmt"
        "net"
        "sync"
        "time"
)

var (
        mutexPlist sync.Mutex
        hchan      chan string
)

func write(conn net.Conn) {
        for {
        	conn.Write([]byte("heartbeat"))
	}
}

const (
        HOST_PEER        = "127.0.0.1"
        PORT_PEER        = "4445"
        HOST_OTHERBN     = "127.0.0.1"
        PORT_OTHERBN     = "4444"
        TIMEOUT_SEC      = 5
        THRESHOLD_SEC    = time.Duration(10) * time.Second
        MAX_PLAYERS_ROOM = 4
)

type Peer struct {
        IPAddress string
        IsLeader  bool
}

func syncToOtherBN() {
        fmt.Println("SYNC")
        //TODO: private connection to other BN

}
func main() {
	conn, err := net.Dial("tcp", "209.131.50.54:4445")
	if err != nil {
		// handle error
	}
        go write(conn)
	dec := gob.NewDecoder(conn)
	var p []Peer
	//p := &Peer{}
	dec.Decode(&p)
	fmt.Println(err)
	fmt.Println(p)
		
}	
