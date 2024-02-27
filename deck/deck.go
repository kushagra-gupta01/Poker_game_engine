package deck

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Spades:
		return "SPADES"
	case Hearts:
		return "HEARTS"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	default:
		panic("Invalid card suit")
	}
}

const (
	Spades   Suit = 0
	Hearts   Suit = 1
	Diamonds Suit = 2
	Clubs    Suit = 3
)

type Card struct {
	suit  Suit
	value int
}

func (c Card) String() string {
	value := strconv.Itoa(c.value)
	if c.value == 1 {
		value = "ACE"
	}
	return fmt.Sprintf("%s of %s %s", value, c.suit, suitToUnicode(c.suit))
}

func NewCard(Suit Suit, Value int) Card {
	if Value > 13 {
		panic("The value of card cannot be greater than 13")
	}
	return Card{
		suit:  Suit,
		value: Value,
	}
}

type Deck [52]Card

func New() Deck {
	var (
		nSuits = 4
		nCards = 13
		d      = [52]Card{}
	)

	x := 0
	for i := 0; i < nSuits; i++ {
		for j := 0; j < nCards; j++ {
			d[x] = NewCard(Suit(i), j+1)
			x++
		}
	}
	return Shuffle(d)
}

func Shuffle(d Deck) Deck {
	for i := 0; i < len(d); i++ {
		r := rand.Intn(i + 1)
		if r != i {
			d[i], d[r] = d[r], d[i]
		}
	}
	return d
}

func suitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("Invalid card suit")
	}
}
