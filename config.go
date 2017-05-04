package main

type Config struct {
	HttpPort int `json:"http_port"`
	P2PPort  int `json:"p2p_port"`
}
