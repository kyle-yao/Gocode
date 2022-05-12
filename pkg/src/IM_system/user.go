package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Ch     chan string
	conn   net.Conn
	server *Server
}

//Create an User Api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		Ch:     make(chan string),
		conn:   conn,
		server: server,
	}
	//Start the listenMessage go routine
	go user.ListenMessage()

	return user
}

//Create a method to inspect current User channel
func (this *User) ListenMessage() {
	for {
		msg := <-this.Ch
		this.conn.Write([]byte(msg + "\n"))
	}
}

//User online
func (this *User) Online() {

	//User online, put them in OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//broadcast User online message
	this.server.BroadCast(this, "is online")
}

//User offline
func (this *User) Offline() {

	//User offline, delete them from OnlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//broadcast User offline message
	this.server.BroadCast(this, "is offline")
}

//Send message to current user
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

//User deal with message
func (this *User) ProcessMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "is online\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {

		//User rename method format
		//rename- yaoxu
		newName := strings.Split(msg, "|")[1]

		//Test whether newName is valid
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("The name exist already, Please try again\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("Success!Your name now is:" + this.Name + "\n")
		}

	} else if len(msg) > 3 && msg[:3] == "to|" {
		//Message format to|username|message

		//Get username and message
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMessage("Invalid format!eg:to|Kyle|Hello!\n")
			return
		} else {
			remoteUser, ok := this.server.OnlineMap[remoteName]
			if !ok {
				this.SendMessage("The user" + remoteName + "not exist!eg:to|Kyle|Hello!\n")
				return
			}
			newMessage := strings.Split(msg, "|")[2]
			if newMessage == "" {
				this.SendMessage("Invalid content!\n")
				return
			}
			remoteUser.SendMessage(this.Name + ":" + newMessage + "\n")
		}

	} else {
		this.server.BroadCast(this, msg)
	}

}
