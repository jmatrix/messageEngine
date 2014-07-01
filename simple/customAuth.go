package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	msgE "github.com/prestonTao/messageEngine"
	"io"
	"net"
	"time"
)

func main() {
	go server()
	time.Sleep(time.Second * 5)
	client()
}

//---------------------------------------------
//          server
//---------------------------------------------
func server() {
	engine := msgE.NewEngine("interServer")
	engine.RegisterMsg(111, RecvMsg)
	engine.SetAuth(new(CustomAuth))
	engine.Listen("127.0.0.1", 9090)
	time.Sleep(time.Second * 10)
}

func RecvMsg(c msgE.Controller, msg msgE.GetPacket) {
	fmt.Println(string(msg.Date))
	session, ok := c.GetSession(msg.Name)
	if ok {
		session.Close()
	}
}

//---------------------------------------------
//          client
//---------------------------------------------
func client() {
	engine := msgE.NewEngine("interClient")
	engine.RegisterMsg(111, RecvMsg)
	engine.SetAuth(new(CustomAuth))
	engine.AddClientConn("test", "127.0.0.1", 9090)

	//给服务器发送消息
	session, _ := engine.GetController().GetSession("test")
	hello := []byte("hello, I'm client")
	session.Send(111, &hello)
	time.Sleep(time.Second * 10)

}

//---------------------------------------------
//          custom Auth
//---------------------------------------------

type CustomAuth struct {
}

//发送
func (this *CustomAuth) SendKey(conn net.Conn, session msgE.Session, name string) (err error) {
	// name := session.GetName()

	lenght := int32(len(name))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, lenght)

	buf.Write([]byte(name))
	conn.Write(buf.Bytes())
	return
}

//接收
func (this *CustomAuth) RecvKey(conn net.Conn) (name string, err error) {
	lenghtByte := make([]byte, 4)
	io.ReadFull(conn, lenghtByte)
	lenght := binary.BigEndian.Uint32(lenghtByte)
	nameByte := make([]byte, lenght)

	n, e := conn.Read(nameByte)
	if e != nil {
		err = e
		return
	}
	name = string(nameByte[:n])

	//检查用户名
	if name == "interClient" {
		//用户验证通过
		fmt.Println("用户验证通过")
		return
	}
	//该用户非法
	err = errors.New("该用户非法")
	fmt.Println("该用户非法")
	return
}
