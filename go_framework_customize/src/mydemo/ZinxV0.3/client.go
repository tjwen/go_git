package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(time.Second)
	//直接链接远程服务器
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	//链接调用Write 写数据

	for {
		_, err := conn.Write([]byte("Hello Zinx v0.3..."))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}
		buf := make([]byte,512)
		cnt, err := conn.Read(buf)
		if err !=nil {
			fmt.Println("read buf error")
		}
		fmt.Printf("read call back : %s , cnt = %d\n", buf,cnt)
	}
	//cpu 阻塞
	time.Sleep(time.Second)

}