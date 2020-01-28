package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

/*
	封包，拆包 模块
	直接面向TCP连接中的数据流，用于TCP 黏包问题
*/

type DataPack struct{}

//拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return  &DataPack{}
}
//获取包头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节)+ID uint32(4字节)
	return 8
}

//封包方法
//dataLen | msgID | data
func (dp *DataPack)Pack (msg ziface.IMessage) ([]byte, error){
	//创建一个存放bytes字节的缓冲
	dataBuffer := bytes.NewBuffer([]byte{})
	//将dataLen写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil{
		return nil,err
	}
	//将MsgID写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil{
		return nil,err
	}
	//将data数据写进dataBuffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil{
		return nil,err
	}
	return dataBuffer.Bytes(),nil
}
//拆包方法 (将包的Head信息读出来） 之后再根据head信息里的data长度， 再进行一次读取
func (dp *DataPack)UnPack(binaryData []byte) (ziface.IMessage,error){
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到dataLen 和 MsgId
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen);err!=nil{
		return nil,err
	}
	//读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id);err!=nil{
		return nil,err
	}
	//判断dataLen是否已经超出了我们允许的最大包长度
	if (utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize){
		return nil,errors.New("too Large msg data recv")
	}
	return msg,nil
}
