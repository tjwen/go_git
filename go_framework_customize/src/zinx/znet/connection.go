package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

/*
	链接模块
*/
type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	//当前链接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	//告知当前链接已经退出的/停止 channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		handleAPI: callbackApi,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			continue
		}
		//调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID", c.ConnID, "handle is error", err)
			break
		}

	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start() ..ConnID:", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前链接写数据的业务
}
func (c *Connection) Stop() {
	fmt.Println("Conn Stop..Conn ID=:", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	close(c.ExitChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) Send(data []byte) error {
	return nil
}
