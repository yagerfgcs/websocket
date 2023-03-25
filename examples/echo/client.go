// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"math/rand"
	"fmt"
)

var addr = flag.String("addr", "localhost:8088", "http service address")
// 定义要发送的消息列表
var messages = []string{"Hello", "World", "WSS", "Message"}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	connect, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer connect.Close()

	// 设置一个不同的随机种子
	rand.Seed(time.Now().UnixNano())

	// 设置一个不同的随机种子
	//randSleep.Seed(time.Now().UnixNano())

	// 执行循环 100 次
	for i := 0; i < 100000000; i++ {
		// 随机选择要发送的消息
		message := messages[rand.Intn(len(messages))]
		// 发送消息
		if err := connect.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 休眠 1 秒钟
		sleepTime := rand.Intn(200)
		time.Sleep(time.Duration(sleepTime)  * time.Millisecond)

		log.Printf("sleepTime:%d", sleepTime)
	}
	// done := make(chan struct{})

	// go func() {
	// 	defer close(done)
	// 	for {
	// 		_, message, err := connect.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			return
	// 		}
	// 		log.Printf("recv: %s", message)
	// 	}
	// }()

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case t := <-ticker.C:
	// 		err := connect.WriteMessage(websocket.TextMessage, []byte(t.String()))
	// 		if err != nil {
	// 			log.Println("write:", err)
	// 			return
	// 		}
	// 	case <-interrupt:
	// 		log.Println("interrupt")

	// 		// Cleanly close the connection by sending a close message and then
	// 		// waiting (with timeout) for the server to close the connection.
	// 		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 		if err != nil {
	// 			log.Println("write close:", err)
	// 			return
	// 		}
	// 		select {
	// 		case <-done:
	// 		case <-time.After(time.Second):
	// 		}
	// 		return
	// 	}
	// }
}
