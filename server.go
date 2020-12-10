package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var listenAddr = "0.0.0.0"
var listenPort = "8888"

//var socksServerIp = "172.18.12.70"
//var socksPort = "10080"
//var mysqlIp = "10.104.20.42"
//var mysqlPort = "3306"

var socksServerIp = "49.234.85.242"
var socksPort = "6666"
var mysqlIp = "127.0.0.1"
var mysqlPort = "3306"

func main() {
	var l net.Listener
	var err error

	l, err = net.Listen("tcp", listenAddr+":"+listenPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + listenAddr + ":" + listenPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		//logs an incoming message
		fmt.Printf("新用户连接上 %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		// Handle connections in a new goroutine.
		go handleServerConn(conn)
	}
}

func handleServerConn(localConn net.Conn) {
	defer localConn.Close()
	var stage = 0

	// 连接远端
	remoteConn, err := net.Dial("tcp", socksServerIp+":"+socksPort)
	if err != nil {
		fmt.Println("远程连接错误:", err.Error())
		return
	}
	defer remoteConn.Close()
	stage = 1

	// 发送socks握手信息
	_, err = remoteConn.Write([]byte{05, 01, 00})
	if err != nil {
		fmt.Println("发送握手失败:", err.Error())
		return
	}
	stage = 2

	// 接收socks握手信息
	buf := make([]byte, 2)
	len, err := remoteConn.Read(buf)
	if err != nil {
		fmt.Println("接收握手失败:", err.Error())
		return
	}
	if bytes.Compare(buf[:len], []byte{05, 00}) != 0 {
		fmt.Println("非标准的socks5握手")
		return
	}
	stage = 3

	// 发送需要连接的地址
	addrBytes := []byte{05, 01, 00}
	ipBytes := Int32ToBytes(StringIpToInt(mysqlIp))
	mysqlPortNumber, _ := strconv.Atoi(mysqlPort)
	portBytes := Int16ToBytes(mysqlPortNumber)
	addrBytes = append(addrBytes, ipBytes...)
	addrBytes = append(addrBytes, portBytes...)
	_, err = remoteConn.Write(addrBytes)
	if err != nil {
		fmt.Println("发送需要连接的地址失败:", err.Error())
		return
	}
	stage = 4

	// 接收socks服务端的远程连接结果
	buf = make([]byte, 1024)
	len, err = remoteConn.Read(buf)
	if err != nil {
		fmt.Println("socks连接数据库失败:", err.Error())
		return
	}
	if bytes.Compare(buf[:2], []byte{05, 00}) != 0 {
		fmt.Println("socks连接数据库失败")
		return
	}
	if len > 10 {
		_, err := localConn.Write(buf[10:len])
		if err != nil {
			fmt.Println("第一次发包失败:", err.Error())
			return
		}
	}
	stage = 5
	fmt.Println(stage)

	//buf = make([]byte, 1024)
	//_, err = remoteConn.Read(buf)
	//fmt.Println(err)
	//fmt.Println(buf)
	var isFinished = false
	wg := sync.WaitGroup{}
	wg.Add(2)
	go handleBindCon(remoteConn, localConn, &wg, &isFinished)
	go handleBindCon(localConn, remoteConn, &wg, &isFinished)

	wg.Wait()
}

func handleBindCon(con1 net.Conn, con2 net.Conn, wg *sync.WaitGroup, isFinished *bool) {
	defer wg.Done()
	defer func() { *isFinished = true }()
	for {
		buf := make([]byte, 1024)
		err := con1.SetDeadline(time.Now().Add(2 * time.Second))
		if err != nil {
			fmt.Println("read时间过长:", err.Error())
			if *isFinished == true {
				return
			}
			continue
		}
		len, err := con1.Read(buf)
		if err != nil {
			if oe, ok := err.(*net.OpError); ok {
				isTimeout := oe.Timeout()
				if isTimeout {
					if *isFinished == true {
						return
					}
					fmt.Println("read时间过长:", err.Error())
					continue
				}
			}
			fmt.Println("收包错误:", err.Error())
			return
		}
		_, err = con2.Write(buf[:len])
		if err != nil {
			fmt.Println("发包错误:", err.Error())
			return
		}
	}
}
