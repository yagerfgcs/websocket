// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/trace"

	"github.com/gorilla/websocket"

	//"runtime/pprof"
	//"runtime"
	"fmt"
	"os/signal"
	"syscall"
)

var addr = flag.String("addr", "localhost:8088", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	fmt.Printf("----start\n")
	// 创建跟踪文件
	f, err := os.Create("./log/trace.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 启动跟踪
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	// // pprof分析
	// pprof_file, err := os.Create("./log/pprof.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer pprof_file.Close()

	// //w, _ := os.OpenFile("pprof.out", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	// //pprof.StartCPUProfile(pprof_file)
	// // 记录 CPU 使用率
	// if err := pprof.StartCPUProfile(pprof_file); err != nil {
	// 	log.Fatal(err)
	// }
	// defer pprof.StopCPUProfile()

	// // 记录内存分配
	// if err := pprof.Lookup("heap").WriteTo(pprof_file, 0); err != nil {
	// 		log.Fatal(err)
	// }

	// // 记录堆栈信息
	// runtime.SetBlockProfileRate(1)
	// defer runtime.SetBlockProfileRate(0)
	// if err := pprof.Lookup("block").WriteTo(pprof_file, 0); err != nil {
	// 		log.Fatal(err)
	// }

	// // 记录锁竞争
	// if err := pprof.Lookup("mutex").WriteTo(pprof_file, 0); err != nil {
	// 		log.Fatal(err)
	// }

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	//log.Fatal(http.ListenAndServe(*addr, nil))
	http.ListenAndServe(*addr, nil)

	// 创建一个信号通道
	sigCh := make(chan os.Signal, 1)

	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	sig := <-sigCh
	fmt.Printf("Received signal %s, shutting down...\n", sig)
	log.Println("----stop CPUProfile")
	//trace.Stop()
	//f.Close()
	// // 关闭文件
	// if err := pprof_file.Close(); err != nil {
	// 	log.Fatal(err)
	// }
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
