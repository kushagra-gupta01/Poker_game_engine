package p2p

type GameStatus uint32

func (g GameStatus) String() string {
	switch g {
	case GameStateWaiting:
		return "WAITING"
	
	case GameStatusDealing:
		return "DEALING"
	
	case GameStatusPreFlop:
		return "PRE-FLOP"

	case GameStatusFlop:
		return "FLOP"

	case GameStatusTurn:
		return "TURN"

	case GameStatusRiver:
		return "RIVER"

	default:
		return "unknown"
	}
}

const (
	GameStateWaiting GameStatus = iota
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type GameState struct {
	isDealer   bool       //should be atomic accessible
	gameStatus GameStatus //should be atomic accessible
}

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) loop() {

}
