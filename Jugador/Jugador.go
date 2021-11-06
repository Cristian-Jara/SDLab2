package main

import (
	"log"
	"fmt"
	"context"
	"time"
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
	
	log.Printf("Response from server: %s", response.Reply)
	if response.Player != "-1"{
		Player = response.Player
		for Status == "Alive" {
			for {
				log.Printf("Presiona 1 para jugar y 2 para ver el monto acumulado actual")
				fmt.Scanln(&input)
				if input == "1" {
					break
				}else if input == "2"{
					response,err := serviceClient.GetMoneyAmount(context.Background(), &pb.MoneyAmount{Money:""})
					if err != nil{
						log.Fatalf("Error when calling GetMoneyAmount: %v",err)
					}
					log.Printf("El monto acumulado es de: "+ response.Money)
				}else{
					log.Printf("Valor erróneo, ingresa nuevamente")
				}
			}
			log.Printf("Esperando a los demás jugadores y al Lider ... ")
			for Playing != true {
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
					}
					log.Printf("Escoge un numero del 1 al 10")
					fmt.Scanln(&Jugada)
					Score += int(Jugada)
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
					ActualRound+=1
					Playing = false
					time.Sleep(1*time.Second)
				}
				if Status == "Alive" && Score >=21 {
					log.Printf("Felicidades, sobreviviste el nivel 1")
				}
				if Score < 21 {
					log.Printf("Mala suerte, no alcanzaste los puntos necesarios, has muerto")
					Status = "Dead"
					break
				}
			} else if ActualStage == "2" {
				ActualRound = 1
				log.Printf("Etapa 2")
				log.Printf("Escoge un numero del 1 al 4")
				fmt.Scanln(&Jugada)
				Score += int(Jugada)
				for Jugada<1 || Jugada >4 {
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
				Playing = false
				time.Sleep(1*time.Second)
				if Status == "Alive" {
					log.Printf("Felicidades, sobreviviste el nivel 2")
				}
			} else if ActualStage == "3" {
				log.Printf("Etapa 3")

				log.Printf("Escoge un numero del 1 al 10")
				fmt.Scanln(&Jugada)
				Score += int(Jugada)
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
				Playing = false
				time.Sleep(1*time.Second)
				if Status == "Alive" {
					log.Printf("Felicidades, sobreviviste el nivel 3, haz ganado!")
					break
				}
			}
		}
	}
}