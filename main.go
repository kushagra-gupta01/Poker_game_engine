package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kushagra-gupta01/Poker_game_engine/deck"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	d:= deck.New()
	fmt.Println(d)

}