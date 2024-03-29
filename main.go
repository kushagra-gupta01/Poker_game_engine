package main

import (
	"log"
	"time"
	// "math/rand"
	// "time"
	// "github.com/kushagra-gupta01/Poker_game_engine/deck"
	"github.com/kushagra-gupta01/Poker_game_engine/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		Version:    "Welcome to Poker V1.1",
		ListenAddr: ":3000",
		GameVariant: p2p.TexasHoldem,
	}

	server := p2p.NewServer(cfg)
	go server.Start()
	time.Sleep(1 * time.Second)

	remoteCfg := p2p.ServerConfig{
		Version:    "Welcome to Poker V1.1",
		ListenAddr: ":4000",
		GameVariant: p2p.TexasHoldem,
	}

	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	if err := remoteServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}
	select {}
	// rand.Seed(time.Now().UnixNano())  ->not needed
	// fmt.Println(deck.New())
}
