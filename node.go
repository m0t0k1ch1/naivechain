package main

import (
	"context"
	"encoding/json"
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
	node.HandleFunc("/blocks", node.blocksHandler)
	node.HandleFunc("/mineBlock", node.mineBlockHandler)

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

func (node *Node) blocksHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(blockchain)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to decode blockchain")
		return
	}

	w.Write(b)
}

func (node *Node) mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to decode data")
		return
	}

	block, err := generateBlock(params.Data)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to generate block")
		return
	}

	if err := blockchain.addBlock(block); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "failed to add block")
		return
	}

	// TODO: bloadcast
}
