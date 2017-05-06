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

type ErrorResponse struct {
	Error string `json:"error"`
}

type Node struct {
	*http.ServeMux
	blockchain *Blockchain
	config     Config
	logger     *log.Logger
}

func newNode(config Config) *Node {
	node := &Node{
		ServeMux:   http.NewServeMux(),
		blockchain: newBlockchain(),
		config:     config,
		logger: log.New(
			os.Stdout,
			"node: ",
			log.Ldate|log.Ltime,
		),
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

func (node *Node) logError(err error) {
	node.logger.Println("[ERROR]", err)
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

func (node *Node) blocksHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(node.blockchain.blocks)
	if err != nil {
		node.error(w, err, "failed to decode blocks")
		return
	}

	node.writeResponse(w, b)
}

func (node *Node) mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		node.error(w, err, "failed to decode data")
		return
	}

	block, err := node.blockchain.generateBlock(params.Data)
	if err != nil {
		node.error(w, err, "failed to generate block")
		return
	}
	blockHash, err := block.hash()
	if err != nil {
		node.error(w, err, "failed to hash block")
		return
	}

	if err := node.blockchain.addBlock(block); err != nil {
		node.error(w, err, "failed to add block")
		return
	}

	// TODO: broadcast

	b, err := json.Marshal(map[string]string{
		"hash": blockHash,
	})

	node.writeResponse(w, b)
}
