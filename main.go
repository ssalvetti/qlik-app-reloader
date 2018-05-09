package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type qlikRPCrequest struct {
	json   string   `json:"jsonrpc"`
	id     int      `json:"id"`
	handle int      `json:"handle"`
	method string   `json:"method"`
	params []string `json:"params"`
}

type qlikRPCresponse struct {
	json   string                            `json:"jsonrpc"`
	id     string                            `json:"id"`
	result map[string]map[string]interface{} `json:"result"`
}

func main() {
	// accept parameters
	app := flag.String("app", "", "app to reload")
	serverAddress := flag.String("server", "localhost:4848", "server host and port, like ip:port")
	flag.Parse()

	// declare where to work
	connStr := fmt.Sprintf("ws://%s/app/%s", *serverAddress, url.QueryEscape(*app))
	c, _, err := websocket.DefaultDialer.Dial(connStr, nil)
	if err != nil {
		log.Fatalf("connection to WS on addr %s failed: %v", connStr, err)
	}

	// prepare main program synchroniaztion
	var wg sync.WaitGroup
	wg.Add(1)

	// setup message parsing
	reply := make(chan []byte, 0)
	done := make(chan struct{}, 0)

	// accept messages from WS conn
	go func() {
		for {
			select {
			case <-done:
				wg.Done()
				return
			case mex := <-reply:
				log.Println("WS message ", mex)
			}
		}
	}()

	// read messages
	go func() {
		defer close(done)
		for {
			_, mex, err := c.ReadMessage()
			if err != nil {
				log.Printf("errore reading from WS conn: %v", err)
				return
			}
			reply <- mex
		}
	}()

	openDoc := &qlikRPCrequest{
		json:   "2.0",
		id:     1,
		method: "OpenDoc",
		handle: -1,
		params: []string{*app},
	}
	if err := c.WriteJSON(openDoc); err != nil {
		log.Fatalf("could not open the document with RPC OpenDoc call: %v", err)
	}

	doReload := &qlikRPCrequest{
		json:   "2.0",
		id:     2,
		method: "DoReload",
		handle: 1,
	}
	if err := c.WriteJSON(doReload); err != nil {
		log.Fatalf("could not reload app with RPC DoReload call: %v", err)
	}
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Fatalf("error clsing WS connection: %v", err)
	}
	wg.Wait()
}
