package main

import (
	"log"
	"fmt"
	"context"
	"time"
	"math/rand"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
)

const (
	LiderIP = "10.6.40.227"
	Puerto = ":50052"
)

var (
	Player = "-1" 
	Status = "Alive"
	Playing = false
	Finished = false
	ActualStage string
	ActualRound int32
	Jugada int32
	Score = 0
	input string
)


func main(){
	conn,err := grpc.Dial(fmt.Sprint(LiderIP,Puerto),grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	serviceClient := pb.NewChatServiceClient(conn)

	message := pb.JoinRequest{ Request: "Play" }

	response,err := serviceClient.JoinToGame(context.Background(),&message)
	if err != nil{
		log.Fatalf("Error when calling SendMessage: %s", err)
	}
	rand.Seed(time.Now().UnixNano())
	log.Printf("Response from server: %s", response.Reply)
	if response.Player != "-1"{
		Player = response.Player
		time.Sleep(5*time.Second)
		for Status == "Alive" {
			log.Printf("Esperando a los demás jugadores y al Lider ... ")
			time.Sleep(3*time.Second)
			for Playing != true {
				time.Sleep(3*time.Second)
				message := pb.GameStarted{Body:Player, Type:""}
				response,err := serviceClient.StageOrRoundStarted(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.Body == "Si"{
					Playing = true
					ActualStage = response.Type
				}else if response.Body == "Killed"{
					Status = "Dead"
					break
				}else if response.Body == "Win"{
					Finished = true
					break
				}
				
			}
			if Status != "Alive"{
				log.Printf("Haz muerto elegido aleatoriamente") 
				break
			}
			if Finished == true {
				log.Printf("Felicidades haz ganado el juego del calamar") 
				break
			}
			Playing = false
			time.Sleep(1*time.Second)
			if ActualStage == "1" {
				ActualRound = 1
				for ActualRound<5 {
					time.Sleep(1*time.Second)
					log.Printf("Esperando inicie la ronda ....")
					for Playing != true {
						message := pb.GameStarted{Body:Player, Type:""}
						response,err := serviceClient.StageOrRoundStarted(context.Background(),&message)
						if err != nil{
							log.Fatalf("Error when calling SendMessage: %s", err)
						}
						if response.Body == "Si"{
							Playing = true
						}
						time.Sleep(1*time.Second)
					}
					time.Sleep(1*time.Second)
					Jugada = int32(rand.Intn(10) + 1) 
					log.Printf("El bot a jugado %v",Jugada)
					Score += int(Jugada)
					message := pb.SendPlay{Player: Player, Plays: Jugada,  Stage: ActualStage, Round: int32(ActualRound), Score: int32(Score)}
					response,err := serviceClient.SendPlays(context.Background(),&message)
					if err != nil{
						log.Fatalf("Error when calling SendMessage: %s", err)
					}
					if response.Alive == false{
						log.Printf("El bot ha muerto")
						Status = "Dead"
						break
					}
					log.Printf("El bot sobrevivió la ronda")
					ActualRound+=1
					Playing = false
					time.Sleep(2*time.Second)
				}
				if Status == "Alive" && Score >=21 {
					log.Printf("Felicidades bot, sobreviviste el nivel 1")
				}
				if Score < 21 {
					log.Printf("Mala suerte, no alcanzaste los puntos necesarios, has muerto")
					Status = "Dead"
					break
				}
			} else if ActualStage == "2" {
				ActualRound = 1
				log.Printf("Etapa 2")
				time.Sleep(2*time.Second)
				Jugada = int32(rand.Intn(4) + 1) 
				log.Printf("El bot ha jugado %v",Jugada)
				Score += int(Jugada)
				message := pb.SendPlay{Player: Player, Plays: Jugada,  Stage: ActualStage, Round: int32(ActualRound), Score: int32(Score)}
				response,err := serviceClient.SendPlays(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.Alive == false{
					log.Printf("El bot ha muerto")
					Status = "Dead"
					break
				}
				Playing = false
				time.Sleep(1*time.Second)
				if Status == "Alive" {
					log.Printf("Felicidades Bot, sobreviviste el nivel 2")
				}
			} else if ActualStage == "3" {
				log.Printf("Etapa 3")
				time.Sleep(5*time.Second)
				Jugada = int32(rand.Intn(10) + 1) 
				log.Printf("El bot ha jugado %v",Jugada)
				message := pb.SendPlay{Player: Player, Plays: Jugada,  Stage: ActualStage, Round: int32(ActualRound), Score: int32(Score)}
				response,err := serviceClient.SendPlays(context.Background(),&message)
				if err != nil{
					log.Fatalf("Error when calling SendMessage: %s", err)
				}
				if response.Alive == false{
					log.Printf("El bot ha muerto")
					Status = "Dead"
				}
				Playing = false
				time.Sleep(1*time.Second)
				if Status == "Alive" {
					log.Printf("Felicidades bot, sobreviviste el nivel 3, haz ganado el Juego!")
				}
			}
		}
	}
}