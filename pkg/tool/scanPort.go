package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var icmp ICMP

//定义ping 输入类型
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

//扫描地址
var ipAddrs chan string = make(chan string)

//扫描结果
var result chan string = make(chan string)

//var presult chan string = make(chan string)

//线程数
var thread chan int = make(chan int)
var nowThread int

//关闭程序
var clo chan bool = make(chan bool)

//定义ip类型
type Ip struct {
	IpList []string `yaml:"ip"`
}

//保存结果
func writeResult() {
	//fileName := "result.txt"
	t1 := time.Now().Format("200601021504")
	fileName := "result_" + t1 + ".log"
	fout, err := os.Create(fileName)
	if err != nil {
		//文件创建失败
		fmt.Println(fileName + " create error")
	}
	defer fout.Close()
	s, ok := <-result
	for ok {
		fout.WriteString(s + "\r\n")
		s, ok = <-result
	}
	//通知进程退出
	clo <- true
}

//根据线程参数启动扫描线程
func runScan() {
	t, ok := <-thread
	nowThread = t
	if ok {
		for i := 0; i < nowThread; i++ {
			go scan(strconv.Itoa(i))
		}
	}
	//等待线程终止
	for <-thread == 0 {
		nowThread--
		if nowThread == 0 {
			//全部线程已终止,关闭结果写入,退出程序
			close(result)
			break
		}
	}
}

//扫描线程
func scan(threadId string) {
	s, ok := <-ipAddrs
	for ok {
		fmt.Println("[thread-" + threadId + "] scan:" + s)
		//fmt.Println("HI", s)
		//_, err := net.Dial("tcp", s)
		//增加超时控制替换原来Dial的调用方法
		_, err := net.DialTimeout("tcp", s, 1*time.Second)
		if err == nil {
			//端口开放
			result <- s
		}

		s, ok = <-ipAddrs
	}
	fmt.Println("[thread-" + threadId + "] end")
	thread <- 0
}

//加载ip配置文件
func OpenFile() []string {
	var ips = make([]string, 0)
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		//return
	}
	var p Ip
	if err = yaml.Unmarshal(bytes, &p); err != nil {
		fmt.Println(err)
		//return
	}

	ips = p.IpList
	return ips

}

//处理参数
func processFlag(arg []string) {

	var ports []int = make([]int, 0)

	tmpPort := os.Args[2]
	if strings.Index(tmpPort, "-") != -1 {
		//连续端口
		tmpPorts := strings.Split(tmpPort, "-")
		var startPort, endPort int
		var err error
		startPort, err = strconv.Atoi(tmpPorts[0])
		if err != nil || startPort < 1 || startPort > 65535 {
			//开始端口不合法
			return
		}
		if len(tmpPorts) >= 2 {
			//指定结束端口
			endPort, err = strconv.Atoi(tmpPorts[1])
			if err != nil || endPort < 1 || endPort > 65535 || endPort < startPort {
				//结束端口不合法
				fmt.Println("'endPort' Setting error")
				return
			}
		} else {
			//未指定结束端口
			endPort = 65535
		}
		for i := 0; startPort+i <= endPort; i++ {
			ports = append(ports, startPort+i)
		}
	} else {
		//一个或多个端口
		ps := strings.Split(tmpPort, ",")
		for i := 0; i < len(ps); i++ {
			p, err := strconv.Atoi(ps[i])
			if err != nil {
				//端口不合法
				fmt.Println("'port' Setting error")
				return
			}
			ports = append(ports, p)

		}
	}
	t, err := strconv.Atoi(os.Args[3])
	//t, err := arg[3]
	if err != nil {
		//线程不合法
		fmt.Println("'thread' Setting error")
		return
	}
	//最大线程5048
	if t < 1 {
		t = 1
	} else if t > 5048 {
		t = 5048
	}

	//传送启动线程数

	thread <- t

	//生成扫描地址列表
	ips := OpenFile()
	il := len(ips)
	for i := 0; i < il; i++ {
		pl := len(ports)
		bool := isping(ips[i])
		if bool == true {
			for j := 0; j < pl; j++ {
				ipAddrs <- ips[i] + ":" + strconv.Itoa(ports[j])
			}
		} else {
			result <- ips[i] + " TimeOut"
		}
	}
	close(ipAddrs)
}

//isping 检查网络函数
func isping(ip string) bool {
	//开始填充数据包
	icmp.Type = 8 //8->echo message  0->reply message
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer

	//先在buffer中写入icmp数据报求去校验和
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	Time, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("ip4:icmp", ip, Time)
	if err != nil {
		return false
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Println("conn.Write error:", err)
		return false
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	num, err := conn.Read(recvBuf)
	if err != nil {
		//log.Println("conn.Read error:", err)
		return false
	}

	conn.SetReadDeadline(time.Time{})

	if string(recvBuf[0:num]) != "" {
		return true
	}
	return false

}

//check  输入数据
func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

//运行程序
func main() {
	flag.Parse()
	if flag.NArg() != 3 && flag.NArg() != 4 {
		//参数不合法
		fmt.Println("正确执行方式，1：IP列表，2：端口范围，3：开启并发数。例如 ./tool ip.yaml 1-65535 1000")
		return
	}
	//获取参数
	args := make([]string, 0, 4)
	for i := 0; i < flag.NArg(); i++ {
		args = append(args, flag.Arg(i))
	}
	//启动扫描线程
	go runScan()
	//启动结果写入线程

	go writeResult()
	//参数处理
	processFlag(args)
	//等待退出指令
	<-clo
	fmt.Println("Exit")
}
