package main

import (
        "encoding/base64"
        "encoding/gob"
        "fmt"
        "io/ioutil"
        "net"
        "net/http"
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

func requestDDNS() {

        client := &http.Client{}

        for {
                user := "jflim91"
                pass := "drawsmore"
                host := "drawsmore.ddns.net"
                ip := HOST_PEER
                format := "http://dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s"
                reqString := fmt.Sprintf(format, user, pass, host, ip)

                req, _ := http.NewRequest("GET", reqString, nil)

                str := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
                fmt.Println(str)
                req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))
                req.Header.Set("User-Agent", "18842 DrawSmore BNode/v0.1")

                resp, err := client.Do(req)
                if err != nil {
                        fmt.Println("DDNS IP update failed")
                } else {
                        defer resp.Body.Close()
                        contents, err := ioutil.ReadAll(resp.Body)
                        if err != nil {
                                fmt.Println("DDNS IP update; unable to retrieve httpBody")
                        }
                        fmt.Println(string(contents))
                        //check response
                        // leave loop if its what we want
                }
                time.Sleep(5 * time.Second)
        }
}

func HeartBeat(conn net.Conn) {

        LastReceived := time.Now()
        for {
                select {
                case m := <-hchan:
                        if m != "" {
                                fmt.Println("Received a heartbeat message", m)
                                LastReceived = time.Now()
                        }

                case <-time.After(TIMEOUT_SEC * time.Second):

                        // check duration since last received heartbeat
                        duration := time.Since(LastReceived)
                        fmt.Println("Duration since last heartbeat ", duration)
                        if duration > THRESHOLD_SEC {
                                fmt.Println("the other BN is gone")
                                fmt.Println("submit self as primary if not already")
                                requestDDNS()
                        }

                        fmt.Println("Sending a heartbeat message")
                        conn.Write([]byte("heartbeat"))
                }
        }
}

func syncToOtherBN() {
        fmt.Println("SYNC")
        //TODO: private connection to other BN

}
func main() {
	conn, err := net.Dial("tcp", "localhost:4445")
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
