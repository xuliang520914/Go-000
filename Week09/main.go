package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)


var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	go sendMessage(conn, user.MessageChannel)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		fmt.Println(strconv.Itoa(user.ID) + ":" + input.Text())
		
		user.MessageChannel <- input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":88")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		
		go handleConn(conn)
	}
}
