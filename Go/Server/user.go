package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn,server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User)Online()  {
	//
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//

	this.server.BroadCast(this, "已上线")
}

func (this *User)Offline(){
		//
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap,this.Name)
		this.server.mapLock.Unlock()
		//
	
		this.server.BroadCast(this, "已下线")
}

func (this *User)SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

func (this *User)Domessage(msg string){

	if msg=="who"{

		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg:="["+user.Addr+"]"+user.Name+":"+"在线....\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

		
	}else if len(msg)>7&&msg[:7]=="rename|"{
		//format: rename|张三
		newName:=strings.Split(msg,"|")[1]

		_,ok:=this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名被使用\n")
		}else{
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName]=this
			this.server.mapLock.Unlock()

			this.Name=newName
			this.SendMsg("您已经更新用户名："+this.Name+"\n")
		}

	}else if(len(msg)>4 && msg[:3]=="to|"){

		//format:to|张三|msg

		//get user name
		remoteName:=strings.Split(msg,"|")[1]
		if remoteName=="" {
			this.SendMsg("消息格式不正确，请使用 \"to|张三|你好啊\"格式。\n")
			return
		}
		//get user object through the user name
		remoteUser,ok:=this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在")
			return	
		}

		//get msg through the user object to send the msg
		content:=strings.Split(msg,"|")[2]
		if content==""{
			this.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name+"对您说："+content)
		
	}else{
		this.server.BroadCast(this,msg)
	}
	
}


func (this *User) ListenMessage() {

	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}

}
