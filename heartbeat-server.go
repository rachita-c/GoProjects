package main

import (
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	HOST_PEER        = "209.131.50.54"
	//PORT_PEER        = "4444"
	HOST_OTHERBN     = "127.0.0.1"
	PORT_OTHERBN     = "4445"
	TIMEOUT_SEC      = 5
	THRESHOLD_SEC    = time.Duration(10) * time.Second
	MAX_PLAYERS_ROOM = 4
)

var (
	mutexPlist sync.Mutex
	hchan      chan string
)

type Peer struct {
	IPAddress string
	IsLeader  bool
}

type Room struct {
	PeerList []Peer // current peers in room
	RoomNo   int
}

var roomList []*Room

func main() {
	fmt.Println("Heerere")
	hchan = make(chan string, 1024)
	fmt.Println("SHA")
	go test()
	syncToOtherBN()
        fmt.Println("SAME:")

}

func test(){
	fmt.Println("CAME")
}
func connectToPeer() {
	//TODO: Wrap it up! Dealing with connection from peers

	/*
	// Listen for incoming connections.
	listenConn, err := net.Listen("tcp", HOST_PEER)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	room := &Room{}
	roomList = append(roomList, room)
	room.RoomNo = 1

	// Close the listener when the application closes.
	defer listenConn.Close()
	Log("BN Listening on " + HOST_PEER + ":" + PORT_PEER)
	for {
		// Listen for an incoming connection.
		if room.isFull() == true {
			newRoom := &Room{}
			newRoom.RoomNo = room.RoomNo + 1
			roomList = append(roomList, newRoom)
			room = roomList[len(roomList)-1]
		}
		conn, err := listenConn.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, room)
	}
	*/
}

// Handles incoming requests.
func handleRequest(conn net.Conn, room *Room) {

	// Send a response back to peer contacting us.
	mutexPlist.Lock()
	encoder := gob.NewEncoder(conn)
	encoder.Encode(room.PeerList)
	Ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	isLeader := (len(room.PeerList) == 0)
	newPeer := Peer{Ip, isLeader}
	room.PeerList = append(room.PeerList, newPeer)
	mutexPlist.Unlock()
	for _, r := range roomList {
		Log(*r)
	}
	// Close the connection when you're done with it.
}

func syncToOtherBN() {
	fmt.Println("SYNC")
	//TODO: private connection to other BN
	ListenToOtherBN()

}

func ListenToOtherBN() {
	fmt.Println("BLAH")
	listenBN, err := net.Listen("tcp", HOST_OTHERBN+":"+PORT_OTHERBN)
	Log("BN Listening on " + HOST_OTHERBN + ":" + PORT_OTHERBN)

	if err != nil {
		fmt.Println("Error listening: ", err.Error())
	}

	defer listenBN.Close()
	for {
		conn, err := listenBN.Accept()
		if err != nil {
			Log("Error Accept: ", err.Error())
			return
		}

		go read(conn)
		go HeartBeat(conn)
	}
}

func read(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			hchan <- ""
			return
		}
		hchan <- string(buffer)
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

func requestDDNS() {

	client := &http.Client{}

	for {
		user := "rchandra101"
		pass := "r1o2k3da"
		host := "bootstrap-rchandra.ddns.net"
		ip := HOST_PEER
		format := "http://dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s"
		reqString := fmt.Sprintf(format, host, ip)

		req, _ := http.NewRequest("GET", reqString, nil)

		str := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
		fmt.Println(str)
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))
		req.Header.Set("User-Agent", "18842 DrawSmore BNode/v0.1")

		resp, err := client.Do(req)
		fmt.Println(resp)
		fmt.Println(err)
		if err != nil {
			Log("DDNS IP update failed")
		} else {
			fmt.Println("YAY")
			defer resp.Body.Close()
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				Log("DDNS IP update; unable to retrieve httpBody")
			}
			fmt.Println(string(contents))
			//check response
			// leave loop if its what we want
		}
		time.Sleep(5 * time.Second)
	}
}

/*
  -----------------------------------------------------------------------------
							Helper Functions
  -----------------------------------------------------------------------------
*/
func Log(v ...interface{}) {
	fmt.Println(v...)
}

func (r Room) isFull() bool {

	if len(r.PeerList) == MAX_PLAYERS_ROOM {
		return true
	}
	return false
}

