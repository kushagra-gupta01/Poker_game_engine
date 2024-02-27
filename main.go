package main

import (
	// "fmt"
	// "math/rand"
	// "time"
	// "github.com/kushagra-gupta01/Poker_game_engine/deck"
	"github.com/kushagra-gupta01/Poker_game_engine/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		ListenAddr: ":3000",
	}

	server := p2p.NewServer(cfg)
	server.Start()
	// rand.Seed(time.Now().UnixNano())  ->not needed
	// d:= deck.New()
	// fmt.Println(d)
}
