package proto

import (
	context "context"
	"fmt"
	"math/rand"
	"strings"
	sync "sync"
	"sync/atomic"
	"time"
)

type GuessNumberBot struct {
	guessNumber int32
	guessChance int
}

const (
	GameName       = "GuessNumber"
	GameDesc       = "Guess A number from 1 to 100 ? You have ten chances to play."
	GuessNumberMax = 100
)

var (
	// game locker , server run one game only
	gameLocker    int32
	sessionLocker = sync.Mutex{}
	guessCounter  = 0
)

func NewGuessNumberBot() *GuessNumberBot {
	guess := rand.New(rand.NewSource(time.Now().Unix())).Int31n(GuessNumberMax) + 1
	fmt.Printf("guess number is %v\n", guess)
	return &GuessNumberBot{
		guessNumber: guess,
		guessChance: 10,
	}
}

// send a request for joining game
func (g *GuessNumberBot) JoinGame(ctx context.Context, jreq *JoinGameRequest) (*JoinGameAnswer, error) {
	sessionLocker.Lock()
	defer sessionLocker.Unlock()

	if strings.Compare(GameName, jreq.GetGameName()) != 0 {
		return nil, fmt.Errorf("we do not run this game %v right now", jreq.GetGameName())
	}

	if atomic.LoadInt32(&gameLocker) != 0 {
		return nil, fmt.Errorf("this game %v is playing by another player", jreq.GetGameName())
	}
	guessCounter = 0
	atomic.StoreInt32(&gameLocker, 1)

	return &JoinGameAnswer{
		ReplyMessage: "you join this game, please give a guess",
		NumMax:       GuessNumberMax,
	}, nil
}

// guess a number ,then server respone a hint
func (g *GuessNumberBot) GuessNumberRight(ctx context.Context, greq *GuessNumber) (*GuessNumberHint, error) {
	if atomic.LoadInt32(&gameLocker) != 1 {
		return nil, fmt.Errorf("you run out of chances")
	}

	if guessCounter >= g.guessChance {
		guessCounter = 0
		atomic.StoreInt32(&gameLocker, 0)
		return nil, fmt.Errorf("you run out of chances")
	}

	guessCounter++
	if greq.GetGuess() < g.guessNumber {
		return &GuessNumberHint{
			Hint:        -1,
			HintMessage: "be less, give an number greater next time",
		}, nil
	} else if greq.GetGuess() > g.guessNumber {
		return &GuessNumberHint{
			Hint:        1,
			HintMessage: "be greater, give a number less next time",
		}, nil
	}
	guessCounter = 0
	atomic.StoreInt32(&gameLocker, 0)
	g.guessNumber = rand.New(rand.NewSource(time.Now().Unix())).Int31n(GuessNumberMax) + 1
	fmt.Printf("guess number is %v\n", g.guessNumber)
	return &GuessNumberHint{
		Hint:        0,
		HintMessage: "bingo",
	}, nil
}

func (g *GuessNumberBot) mustEmbedUnimplementedGuessNumberGameServer() {}
