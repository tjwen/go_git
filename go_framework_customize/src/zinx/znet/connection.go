package znet

import (
	"errors"
	"fmt"
	"io"
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
	//handleAPI ziface.HandleFunc
	//告知当前链接已经退出的/停止 channel (由Reader告知Writer退出）
	ExitChan chan bool
	//无缓冲读管道，用于读，写Goroutine 之间的消息通信
	msgChan chan []byte

	//该链接处理的方法Router
	//Router ziface.IRouter

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32 /*, router ziface.IRouter*/, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false, //isClosed为true 意思是关闭
		msgChan:  make(chan []byte),
		//handleAPI: callbackApi,  //不在需要了 换成Router
		//Router: router,
		MsgHandle: msgHandler,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

/*
	写消息Goroutine ,专门发送给客户端写消息的模块
*/

func (c * Connection) StartWriter()  {
	fmt.Println("[Writer Goroutine ns running]")
	defer fmt.Println(c.RemoteAddr().String(),"[conn Writer Exit!]")
	//不断阻塞等待channel的消息，进行写给客户端
	for{
		select {
			case data := <-c.msgChan:
				//有数据要写给客户端
				if _,err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error", err)
					return
				}
		case <-c.ExitChan:
			//代表Reader
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		/*//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			continue
		}*/
		/*
			//调用当前链接所绑定的HandleAPI
			if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
				fmt.Println("ConnID", c.ConnID, "handle is error", err)
				break
			}
		*/
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 的二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpacke error", err)
			break
		}
		//拆包，得到msgId 和 msgDatalen 放在msg消息中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error :", err)
				break
			}
		}
		msg.SetData(data)

		//根据datalen 再次读取Data，放在msg.Data中

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//从路由中，找到注册绑定的Conn对应得router调用
		/*go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)*/
		//根据绑定好的MsgID 找到对应处理api业务 执行
		go c.MsgHandle.DoMsgHandler(&req)
	}
}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ..ConnID:", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前链接写数据的业务
	//启动从当前链接写数据的业务
	go c.StartWriter()

}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop..Conn ID=:", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket链接
	c.Conn.Close()
	//告知Writer关闭
	c.ExitChan <-true
	//关闭管道 回收资源
	close(c.ExitChan)
	close(c.msgChan)

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

//提供一个sendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDtaLen/MsgId/Data
	dp := NewDataPack()

	//MsgDataLen/MsgID/Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgId)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	/*if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id", msgId, "error:", err)
		return errors.New("conn Write error")
	}*/
	c.msgChan <-binaryMsg
	return nil
}
