package znet

import "zinx/ziface"

type  Request struct {
	//已经和客户建立好的链接
	conn ziface.IConnection
	//客户端请求的数据
	//data []byte
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//func (r *Request) GetData() []byte {
//	return r.data()
//}

//ZinxV0.5中修改 data ->msg
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
//ZinxV0.5中增加的
func(r *Request) GetMsgID() uint32{
	return r.msg.GetMsgId()
}