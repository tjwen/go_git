package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(1*time.Second)
	//直接链接远程服务器
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	//链接调用Write 写数据

	for {
		//发送封包的mssage消息 MsgID:0 
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxVO.6 client Test Message")))
		if err != nil{
			fmt.Println("Pack error:",err)
			return
		}
		if _, err := conn.Write(binaryMsg);err != nil{
			fmt.Println("write error",err)
			return
		}

		//服务器就应该给我们回复一个message数据，MsgID:1 ping...ping...ping

		binaryHead := make([]byte,dp.GetHeadLen())
		if _,err:=io.ReadFull(conn,binaryHead); err !=nil {
			fmt.Println("read head error",err)
			break
		}

		//先读取流中的head部分 得到ID 和 dataLen
		msgHead, err := dp.UnPack(binaryHead)
		if err!= nil{
			fmt.Println("clent unpack msgHead error",err)
			break
		}

		if msgHead.GetMsgLen()>0 {
			//msg里是有数据的，
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte,msg.GetMsgLen())

			if _,err:= io.ReadFull(conn,msg.Data); err != nil{
				fmt.Println("read msg data error,",err)
				return
			}
			fmt.Println("----->Recv Server Msg ID=",msg.Id,"len=",msg.DataLen,"data=",string(msg.Data))
		}

		//再根据DataLen进行第二次读取，将data读出来



		time.Sleep(1*time.Second)

	}
	//cpu 阻塞

}