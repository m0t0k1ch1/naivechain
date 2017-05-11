package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/net/websocket"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Node struct {
	*http.ServeMux
	blockchain *Blockchain
	conns      []*Conn
	mu         sync.RWMutex
	logger     *log.Logger
}

func newNode() *Node {
	return &Node{
		blockchain: newBlockchain(),
		conns:      []*Conn{},
		mu:         sync.RWMutex{},
		logger: log.New(
			os.Stdout,
			"node: ",
			log.Ldate|log.Ltime,
		),
	}
}

func (node *Node) newApiServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/blocks", node.blocksHandler)
	mux.HandleFunc("/mineBlock", node.mineBlockHandler)
	mux.HandleFunc("/peers", node.peersHandler)
	mux.HandleFunc("/addPeer", node.addPeerHandler)

	return &http.Server{
		Handler: mux,
		Addr:    *apiAddr,
	}
}

func (node *Node) newP2PServer() *http.Server {
	return &http.Server{
		Handler: websocket.Handler(func(ws *websocket.Conn) {
			conn := newConn(ws)
			node.log("connect to peer:", conn.remoteHost())
			node.addConn(conn)
			node.p2pHandler(conn)
		}),
		Addr: *p2pAddr,
	}
}

func (node *Node) run() {
	apiSrv := node.newApiServer()
	go func() {
		node.log("start HTTP server for API")
		if err := apiSrv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	p2pSrv := node.newP2PServer()
	go func() {
		node.log("start WebSocket server for P2P")
		if err := p2pSrv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM)
	for {
		s := <-signalCh
		if s == syscall.SIGTERM {
			node.log("stop servers")
			apiSrv.Shutdown(context.Background())
			p2pSrv.Shutdown(context.Background())
		}
	}
}

func (node *Node) log(v ...interface{}) {
	node.logger.Println(v)
}

func (node *Node) logError(err error) {
	node.log("[ERROR]", err)
}

func (node *Node) writeResponse(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (node *Node) error(w http.ResponseWriter, err error, message string) {
	node.logError(err)

	b, err := json.Marshal(&ErrorResponse{
		Error: message,
	})
	if err != nil {
		node.logError(err)
	}

	node.writeResponse(w, b)
}
