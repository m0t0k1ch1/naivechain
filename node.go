package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Node struct {
	*http.ServeMux
	config Config
}

func newNode(config Config) *Node {
	node := &Node{
		ServeMux: http.NewServeMux(),
		config:   config,
	}
	node.HandleFunc("/blocks", blocksHandler)
	node.HandleFunc("/mineBlock", mineBlockHandler)

	return node
}

func (node *Node) run() {
	httpSrv := &http.Server{
		Handler: node,
		Addr:    fmt.Sprintf(":%d", node.config.HttpPort),
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM)
	for {
		s := <-signalCh
		if s == syscall.SIGTERM {
			httpSrv.Shutdown(context.Background())
		}
	}
}
