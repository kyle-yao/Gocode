package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//online usermap
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//Message channel
	Message chan string
}

// Create a server port
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//Method to broadcast Message
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

//Inspect BroadCast channel
func (this *Server) Listener() {
	for {
		msg := <-this.Message

		//Send Message to All online user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.Ch <- msg
		}

		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	//
	// fmt.Println("Linked success")
	user := NewUser(conn, this)
	//User get online, Add user to onlineMap
	user.Online()

	isLive := make(chan bool)
	//receive message from user
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error:", err)
				return
			}

			//Extract message(Drop '\n')
			msg := string(buf[:n-1])
			//User deal with message
			user.ProcessMessage(msg)

			isLive <- true
		}
	}()
	//current channel block
	for {
		select {
		case <-isLive:
			//Do nothing to activate the time.After function
		case <-time.After(time.Second * 60):
			//Exceed time limit, then kick user offline

			user.SendMessage("Time out! You are kicked off!")
			//Close user channel
			close(user.Ch)
			conn.Close()

			//Exit current Handler
			//runtime.Goexit()
			return
		}
	}

}

// Interface to start server
func (this *Server) Start() {
	//socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	//close listen socket
	defer listen.Close()
	//Start listenMessage go routine
	go this.Listener()

	//accept
	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			fmt.Println("Accept error:", err2)
			continue
		}

		go this.Handler(conn)
	}
	//do handler

}
