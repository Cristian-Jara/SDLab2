package main

import (
	"log"
	"fmt"
	"context"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
)

const (
	LiderIP = "10.6.40.227"
	LocalIP = "localhost"
	Puerto = ":50052"
)

var (
	Player = "-1" 
	Status = "Alive"
	Playing = false
	ActualStage string
	ActualRound int32
	Jugada int32
	Score = 0
)

func main(){
	conn,err := grpc.Dial(fmt.Sprint(LocalIP,Puerto),grpc.WithInsecure())
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
		Player = response.Player
		for Status == "Alive" {
			for Playing != true {
				message := pb.GameStarted{Body:"", Type:""}
				response,err := serviceClient.StageOrRoundStarted(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.Body == "Si"{
					Playing = true
					ActualStage = response.Type
				} 
			}
			if ActualStage == "1" {
				log.Printf("Escoge un numero del 1 al 10")
				fmt.Scanln(&Jugada)
				for Jugada<1 || Jugada >10 {
					log.Printf("Escoge un numero válido")
					fmt.Scanln(&Jugada)
				}
				message := pb.SendPlay{Player: Player, Plays: Jugada,  Stage: ActualStage, Round: int32(ActualRound), Score: int32(Score)}
				response,err := serviceClient.SendPlays(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.Alive == false{
					log.Printf("Haz muerto")
					Status = "Dead"
					break
				}
				log.Printf("Sobreviviste la ronda")
			}
		}
	}
}