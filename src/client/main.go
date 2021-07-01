package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/cloudbiter/guess-number-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGuessNumberGameClient(conn)

	joinGameReq := pb.JoinGameRequest{
		GameName: "GuessNumber",
	}

	joinGameAns, err := client.JoinGame(context.Background(), &joinGameReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("reply from game server %v, you guess a number less thant %v\n", joinGameAns.GetReplyMessage(), joinGameAns.GetNumMax())
	begin := int32(1)
	end := joinGameAns.GetNumMax()

	var guess int
	fmt.Scanf("%d", &guess)

	gn := pb.GuessNumber{
		Guess: int32(guess),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for {
		gnh, err := client.GuessNumberRight(ctx, &gn)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("guess: %v begin: %v end: %v hint: %v hintmessage: %v\n", gn.GetGuess(), begin, end, gnh.GetHint(), gnh.GetHintMessage())
		fmt.Printf("guess: %v hint: %v hintmessage: %v\n", gn.GetGuess(), gnh.GetHint(), gnh.GetHintMessage())
		hint := gnh.GetHint()
		if hint == 0 {
			break
		} else if hint == -1 {
			//end = gn.GetGuess() - begin
			begin = gn.GetGuess()
			nextGuess := guessForMe(begin, end)
			gn = pb.GuessNumber{
				Guess: nextGuess,
			}
		} else if hint == 1 {
			end = gn.GetGuess()
			nextGuess := guessForMe(begin, end)
			gn = pb.GuessNumber{
				Guess: nextGuess,
			}
		}
	}
}

func guessForMe(begin, end int32) int32 {
	//return begin + rand.New(rand.NewSource(time.Now().Unix())).Int31n(end-begin)
	return (end + begin) / 2
}
