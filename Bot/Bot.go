package main

import (
	"log"
	"fmt"
	"math/rand"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
)

const (
	LiderIP2 = "10.6.40.227"
	LocalIP2 = "localhost"
	Puerto2 = ":50052"
)

var (
	Bot = "-1" 
	BotStatus = "Alive"
	BotPlaying = false
	BotActualStage = ""
	BotJugada int
)

func main(){
	conn,err := grpc.Dial(fmt.Sprint(LocalIP2,Puerto2),grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	serviceClient := pb.NewChatServiceClient(conn)

	message := pb.JoinRequest{ Request: "Play" }

	response,err := serviceClient.JoinToGame(context.Background(),&message)
	if err != nil{
		log.Fatalf("Error when calling SendMessage: %s", err)
	}
	
	log.Printf("Response from server: %s", response.Reply)
	if response.Player != "-1"{
		Bot = response.Player
		for BotStatus == "Alive" {
			for BotPlaying != true {
				message := pb.GameStarted{Body:"", Type:""}
				response,err := serviceClient.StageOrRoundStarted(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.Body == "Si"{
					BotPlaying = true
					BotActualStage = response.Type
				} 
			}
			if BotActualStage == "Stage 1" {
				log.Printf("Escoge un numero del 1 al 10")
				BotJugada = rand.Intn(10) + 1
				message := pb.Play{Plays: int32(BotJugada), Player: Bot}
				response,err := serviceClient.PlayTheGame(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.RoundResult == "Moriste"{
					log.Printf("Haz muerto")
					BotStatus = "Dead"
					break
				}
				log.Printf("Sobreviviste la ronda")

			}

		}
	}
}