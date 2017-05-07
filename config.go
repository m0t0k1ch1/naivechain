package main

type Config struct {
	Api *ApiConfig `json:"api"`
	P2P *P2PConfig `json:"p2p"`
}

type ApiConfig struct {
	Port int `json:"port"`
}

type P2PConfig struct {
	Port   int    `json:"port"`
	Origin string `json:"origin"`
}
