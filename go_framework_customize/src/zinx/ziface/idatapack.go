package ziface

/*
	封包，拆包 模块
	直接面向TCP连接中的数据流，用于TCP 黏包问题
*/

type IDataPack interface {
	//获取包头的长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(msg IMessage) ([]byte,error)
	//拆包方法
	UnPack([]byte) (IMessage ,error)
}

