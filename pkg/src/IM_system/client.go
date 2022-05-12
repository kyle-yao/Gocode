package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClinet(serverIP string, serverPort int) *Client {
	//Create Client Object
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       999,
	}
	//Link to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn
	//Return Object
	return client
}

//deal with response from server
func (client *Client) DealResponse() {
	//once client.conn get data, copy to stdout output. Permanent block
	io.Copy(os.Stdout, client.conn)
	/*for {
		buf := make()
		client.conn.Read(buf)
		fmt.Println(buf)
	}
	*/
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.PUBLIC CHAT")
	fmt.Println("2.PRIVATE CHAT")
	fmt.Println("3.RENAME")
	fmt.Println("0.Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>Please enter valid number<<<<")
		return false
	}

}

//flag=1 call PublicChat()
func (client *Client) PublicChat() {

	var chatMsg string
	fmt.Println("Please Enter message.Input\"q\" to exit.")
	fmt.Scanln(&chatMsg)

	for chatMsg != "q" {

		//send to server
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write error:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("Please Enter message.Input\"q\" to exit.")
		fmt.Scanln(&chatMsg)
	}
}

//flag=2 call PrivateChat()
func (clinet *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := clinet.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write error:", err)
		return
	}

}

func (client *Client) PrivateChat() {

	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>Please enter frined name, enter \"q\" to exit:")
	fmt.Scanln(&remoteName)

	for remoteName != "q" {
		fmt.Println("Please enter content, q exit:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "q" {
			if len(chatMsg) != 0 {

				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write error:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("Please enter content, q exit:")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>Please enter frined name, enter \"q\" to exit:")
		fmt.Scanln(&remoteName)
	}

}

//flag=3 call UpdateName()
func (client *Client) UpdateName() bool {
	fmt.Println("Please enter username:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		//Switch to different mode according to flag
		switch client.flag {
		case 1:
			//public chat
			fmt.Println("Choose public chat")
			client.PublicChat()
			break
		case 2:
			//private chat
			fmt.Println("Choose private chat")
			client.PrivateChat()
			break
		case 3:
			//rename mode
			client.UpdateName()
			break

		}
	}
}

var serverPort int
var serverIP string

//./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "set serverIp(Default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set serverPort(Default: 8888)")
}

func main() {
	//command parsing
	flag.Parse()

	client := NewClinet(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>>Linked Failed.")
		return
	}
	//Create a gorountine to deal with response from server
	go client.DealResponse()

	fmt.Println(">>>>Linked Success!")

	//Start Client
	client.Run()
}
