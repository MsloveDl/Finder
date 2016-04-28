// Finder project main.go
package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Loop(startPort, endPort int, inPort chan int) {
	for i := startPort; i <= endPort; i++ {
		inPort <- i
	}
}

func Scanner(inPort, outPort, out chan int, ip net.IP, endPort int) {
	for {
		in := <-inPort
		tcpAddr := &net.TCPAddr{IP: ip, Port: in}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if nil != err {
			outPort <- 0
		} else {
			outPort <- in
			conn.Close()
		}

		if in == endPort {
			out <- in
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	inPort := make(chan int)
	startTime := time.Now().Unix()
	outPort := make(chan int)
	out := make(chan int)
	collect := []int{}

	if 4 != len(os.Args) {
		fmt.Println("Usage: Finder.exe IP startPort endPort")
		fmt.Println("Endport must be larger than startPort")

		os.Exit(0)
	}

	ip := net.ParseIP(os.Args[1])

	if os.Args[3] < os.Args[2] {
		fmt.Println("Usage: Finder.exe IP startPort endPort")
		fmt.Println("Endport must be larger than startPort")

		os.Exit(0)
	}

	fmt.Printf("The ip is %s \r\n", ip)

	startPort, _ := strconv.Atoi(os.Args[2])
	endPort, _ := strconv.Atoi(os.Args[3])

	fmt.Printf("%d------%d \r\n", startPort, endPort)

	// 启动扫描控制线程
	go Loop(startPort, endPort, inPort)

	// 等待管道返回的信息
	for {
		select {
		case <-out:
			fmt.Println(collect)
			endTime := time.Now().Unix()

			fmt.Println("The scan process has spent", endTime-startTime, "second")

			os.Exit(0)

		default:
			// 启动扫描作业线程
			go Scanner(inPort, outPort, out, ip, endPort)
			port := <-outPort

			if 0 != port {
				collect = append(collect, port)
			}
		}
	}
}
