package main

import (
	"context"
	"fmt"

	pb "github.com/fabiusinfo/SquidGame/proto"
	"google.golang.org/grpc"

	"log"
	"math/rand"
	"strconv"
	"time"
)

type PlayerStruct struct {
	id    string
	alive bool
	round int32
	score int32
}

func main() {
	var action string
	playerNumber := "1"
	play := "2"
	actualStage := "none"
	codes1 := "none"
	codes2 := "none"
	codes3 := "none"
	next:="none"
	started := false
	flag1 := false

	var list_of_players []PlayerStruct

	//inscripción
	for !flag1 {
		fmt.Println("escribe join para inscribirse en el SquidGame: ")
		fmt.Scanln(&action)
		if action == "join" {
			flag1 = true
		}
	}

	if actualStage == "none" {
		conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())
		if err != nil {
			panic("cannot connect with server " + err.Error())
		}
		servicePlayer := pb.NewSquidGameServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := servicePlayer.JoinGame(ctx, &pb.JoinRequest{Player: playerNumber})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		list_of_players = append(list_of_players, PlayerStruct{playerNumber, true, 1, 0})
		log.Printf("inscrito")
		codes1 = r.GetCodes1()
		codes2 = r.GetCodes2()
		codes3 = r.GetCodes3()
		actualStage = codes1
		started = true
		
		if started == true {
			fmt.Println("inscripción al SquidGame realizada con éxito.")
		}
	} else {
		fmt.Println("ya estas inscrito.")
	}

	//inscribimos los bots

	for i := 1; i < 16; i++ {

		list_of_players = append(list_of_players, PlayerStruct{strconv.Itoa(i + 1), true, 1, 0})

		conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

		if err != nil {
			panic("cannot connect with server " + err.Error())
		}

		servicePlayer := pb.NewSquidGameServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := servicePlayer.JoinGame(ctx, &pb.JoinRequest{Player: list_of_players[i].id})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("inscrito")

		list_of_players[i].id = strconv.Itoa(i + 1)
		list_of_players[i].alive = true

		codes2 = r.GetCodes2()

		fmt.Println("inscripción al SquidGame realizada con éxito.")
	}
	fmt.Println(codes1 + codes2 + codes3)
	//Aquí finaliza la inscripción


	started=false
	flag1 = false
	for !flag1 {
		fmt.Println("ingresa next para comenzar el nivel 1")
		fmt.Scanln(&next)
		if next == "next" {
			conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

		if err != nil {
			panic("cannot connect with server " + err.Error())
		}

		servicePlayer := pb.NewSquidGameServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := servicePlayer.Started(ctx, &pb.StartRequest{Message: "solicito ingresar al nivel 1"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		started=r.GetStarted()
			if started==true{
				flag1 = true
				break
			}
			
			
		}else{
			fmt.Println("ingresaste mal el comando")
		}
	}

	

	//Aquí realizar jugada o checkAmount nivel 1
	contStage := 0
	actualStage = "1rv"
	fmt.Println(actualStage)
	flag1 = false
	for !flag1 {
		if contStage == 4 {
			flag1 = true
			break
		}
		fmt.Println("STAGE 1: escribe send -> enviar jugada, check -> solicitar monto: ")
		fmt.Scanln(&action)
		if action == "send" {
			contStage += 1
			if list_of_players[0].alive == true {
				if list_of_players[0].score < 21 {
					fmt.Println("escribe un número del 1 al 10: ")
					fmt.Scanln(&play)

					conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

					if err != nil {
						panic("cannot connect with server " + err.Error())
					}

					servicePlayer := pb.NewSquidGameServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: playerNumber, Play: play, Stage: actualStage, Round: list_of_players[0].round, Score: list_of_players[0].score})
					if err != nil {
						log.Fatalf("fallo 1: %v", err)
					}
					actualStage = r.GetStage()
					list_of_players[0].round = r.GetRound()
					list_of_players[0].alive = r.GetAlive() // el jugador debe estar en la posicion 0 de la lista
					play_int, err32 := strconv.Atoi(play)
					if err32 != nil {
						log.Fatalf("fallo 32: %v", err32)
					}
					list_of_players[0].score = list_of_players[0].score + int32(play_int)
				} else {
					fmt.Println(" lograste sumar 21, estas salvado")
				}
			} else {
				fmt.Println("el jugador está muerto")
			}
			// sección bots

			//Este pedazo de código es para las jugadas de los bots
			for i := 1; i < 16; i++ {

				if list_of_players[i].alive == true {
					if list_of_players[i].score < 21 {
						fmt.Println(strconv.Itoa(i))
						botPlayer := list_of_players[i].id
						conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

						if err != nil {
							panic("cannot connect with server " + err.Error())
						}

						servicePlayer := pb.NewSquidGameServiceClient(conn)

						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						//jugada aleatoria
						rand.Seed(time.Now().UnixNano())
						playsend := rand.Int63n(10) + 1
						playsend_str := strconv.Itoa(int(playsend))
						r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: botPlayer, Play: playsend_str, Stage: actualStage, Round: list_of_players[i].round})
						if err != nil {
							log.Fatalf("fallo 1: %v", err)
						}
						actualStage = r.GetStage()
						list_of_players[i].round = r.GetRound()
						list_of_players[i].alive = r.GetAlive()
						list_of_players[i].score = list_of_players[i].score + int32(playsend)
					} else {
						fmt.Println(" lograste sumar 21, estas salvado")
					}
				}
				fmt.Println(list_of_players[i])
			}

		} else if action == "check" {
			message := "solicito monto"

			conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

			if err != nil {
				panic("cannot connect with server " + err.Error())
			}

			servicePlayer := pb.NewSquidGameServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := servicePlayer.AmountCheck(ctx, &pb.AmountRequest{Message: message})
			if err != nil {
				log.Fatalf("no se pudo solicitar el monto: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMonto())

		} else {
			fmt.Println("ingresaste mal el comando")
		}
	}
		
		started=false
		flag1 = false
		for !flag1 {
			fmt.Println("ingresa next para comenzar el nivel 2")
			fmt.Scanln(&next)
			if next == "next" {
				conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())
	
			if err != nil {
				panic("cannot connect with server " + err.Error())
			}
	
			servicePlayer := pb.NewSquidGameServiceClient(conn)
	
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := servicePlayer.Started(ctx, &pb.StartRequest{Message: "solicito ingresar al nivel 2"})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			started=r.GetStarted()
				if started==true{
					flag1 = true
					break
				}
				
				
			}else{
				fmt.Println("ingresaste mal el comando")
			}
		}
	
		fmt.Println(actualStage)
		for i := 0; i < 16; i++ {
			list_of_players[i].score = 0
			if list_of_players[i].alive == true {
				conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

						if err != nil {
							panic("cannot connect with server " + err.Error())
						}

						servicePlayer := pb.NewSquidGameServiceClient(conn)

						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()

						r, err := servicePlayer.DeadOrAlive(ctx, &pb.DeadRequest{Player:list_of_players[i].id , Stage: actualStage })
						if err != nil {
							log.Fatalf("fallo 1: %v", err)
						}
						list_of_players[i].alive=r.GetDead()
			}
		}



		//Aquí realizar jugada o checkAmount nivel 2 
	contStage = 0
	actualStage = "2tc"
	fmt.Println(actualStage)
	flag1 = false
	for !flag1 {
		if 1 >= len(list_of_players) {
			flag1 = false
			break
		}
		if contStage == 1 {
			flag1 = true
			break
		}
		fmt.Println("STAGE 2: escribe send -> enviar jugada, check -> solicitar monto: ")
		fmt.Scanln(&action)
		if action == "send" {
			contStage += 1
			if list_of_players[0].alive == true {
				
					fmt.Println("escribe un número del 1 al 4: ")
					fmt.Scanln(&play)

					conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

					if err != nil {
						panic("cannot connect with server " + err.Error())
					}

					servicePlayer := pb.NewSquidGameServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: playerNumber, Play: play, Stage: actualStage, Round: list_of_players[0].round, Score: list_of_players[0].score})
					if err != nil {
						log.Fatalf("fallo 1: %v", err)
					}
					actualStage = r.GetStage()
					list_of_players[0].round = r.GetRound()
					list_of_players[0].alive = r.GetAlive() // el jugador debe estar en la posicion 0 de la lista
					play_int, err32 := strconv.Atoi(play)
					if err32 != nil {
						log.Fatalf("fallo 32: %v", err32)
					}
					list_of_players[0].score = list_of_players[0].score + int32(play_int)
				
			} else {
				fmt.Println("el jugador está muerto")
			}
			// sección bots

			//Este pedazo de código es para las jugadas de los bots
			for i := 1; i < 16; i++ {

				if list_of_players[i].alive == true {
						fmt.Println(strconv.Itoa(i))
						botPlayer := list_of_players[i].id
						conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

						if err != nil {
							panic("cannot connect with server " + err.Error())
						}

						servicePlayer := pb.NewSquidGameServiceClient(conn)

						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						//jugada aleatoria
						rand.Seed(time.Now().UnixNano())
						playsend := rand.Int63n(3) + 1
						playsend_str := strconv.Itoa(int(playsend))
						r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: botPlayer, Play: playsend_str, Stage: actualStage, Round: list_of_players[i].round})
						if err != nil {
							log.Fatalf("fallo 1: %v", err)
						}
						actualStage = r.GetStage()
						list_of_players[i].round = r.GetRound()
						list_of_players[i].alive = r.GetAlive()
						list_of_players[i].score = list_of_players[i].score + int32(playsend)
					
				}
				fmt.Println(list_of_players[i])
			}

		} else if action == "check" {
			message := "solicito monto"

			conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

			if err != nil {
				panic("cannot connect with server " + err.Error())
			}

			servicePlayer := pb.NewSquidGameServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := servicePlayer.AmountCheck(ctx, &pb.AmountRequest{Message: message})
			if err != nil {
				log.Fatalf("no se pudo solicitar el monto: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMonto())

		} else {
			fmt.Println("ingresaste mal el comando")
		}
	}
	
		//Aquí realizar jugada o checkAmount nivel 2 
	contStage = 0
	actualStage = "2tc"
	fmt.Println(actualStage)
	flag1 = false
	for !flag1 {
		if 1 >= len(list_of_players) {
			flag1 = false
			break
		}
		if contStage == 1 {
			flag1 = true
			break
		}
		fmt.Println("STAGE 2: escribe send -> enviar jugada, check -> solicitar monto: ")
		fmt.Scanln(&action)
		if action == "send" {
			contStage += 1
			if list_of_players[0].alive == true {
				
					fmt.Println("escribe un número del 1 al 4: ")
					fmt.Scanln(&play)

					conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

					if err != nil {
						panic("cannot connect with server " + err.Error())
					}

					servicePlayer := pb.NewSquidGameServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: playerNumber, Play: play, Stage: actualStage, Round: list_of_players[0].round, Score: list_of_players[0].score})
					if err != nil {
						log.Fatalf("fallo 1: %v", err)
					}
					actualStage = r.GetStage()
					list_of_players[0].round = r.GetRound()
					list_of_players[0].alive = r.GetAlive() // el jugador debe estar en la posicion 0 de la lista
					play_int, err32 := strconv.Atoi(play)
					if err32 != nil {
						log.Fatalf("fallo 32: %v", err32)
					}
					list_of_players[0].score = list_of_players[0].score + int32(play_int)
				
			} else {
				fmt.Println("el jugador está muerto")
			}
			// sección bots

			//Este pedazo de código es para las jugadas de los bots
			for i := 1; i < 16; i++ {

				if list_of_players[i].alive == true {
						fmt.Println(strconv.Itoa(i))
						botPlayer := list_of_players[i].id
						conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

						if err != nil {
							panic("cannot connect with server " + err.Error())
						}

						servicePlayer := pb.NewSquidGameServiceClient(conn)

						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						//jugada aleatoria
						rand.Seed(time.Now().UnixNano())
						playsend := rand.Int63n(3) + 1
						playsend_str := strconv.Itoa(int(playsend))
						r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: botPlayer, Play: playsend_str, Stage: actualStage, Round: list_of_players[i].round})
						if err != nil {
							log.Fatalf("fallo 1: %v", err)
						}
						actualStage = r.GetStage()
						list_of_players[i].round = r.GetRound()
						list_of_players[i].alive = r.GetAlive()
						list_of_players[i].score = list_of_players[i].score + int32(playsend)
					
				}
				fmt.Println(list_of_players[i])
			}

		} else if action == "check" {
			message := "solicito monto"

			conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

			if err != nil {
				panic("cannot connect with server " + err.Error())
			}

			servicePlayer := pb.NewSquidGameServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := servicePlayer.AmountCheck(ctx, &pb.AmountRequest{Message: message})
			if err != nil {
				log.Fatalf("no se pudo solicitar el monto: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMonto())

		} else {
			fmt.Println("ingresaste mal el comando")
		}
	}



	//Aquí realizar jugada o checkAmount nivel 3
	contStage = 0
	actualStage = "3tn"
	fmt.Println(actualStage)
	for !flag1 {
		if 1 >= len(list_of_players) {
			flag1 = false
			break
		}

		if contStage == 1 {
			flag1 = true
			break
		}
		fmt.Println("STAGE 3: escribe send -> enviar jugada, check -> solicitar monto: ")
		fmt.Scanln(&action)
		if action == "send" {
			contStage += 1
			if list_of_players[0].alive == true {
				
					fmt.Println("escribe un número del 1 al 10: ")
					fmt.Scanln(&play)

					conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

					if err != nil {
						panic("cannot connect with server " + err.Error())
					}

					servicePlayer := pb.NewSquidGameServiceClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: playerNumber, Play: play, Stage: actualStage, Round: list_of_players[0].round, Score: list_of_players[0].score})
					if err != nil {
						log.Fatalf("fallo 1: %v", err)
					}
					actualStage = r.GetStage()
					list_of_players[0].round = r.GetRound()
					list_of_players[0].alive = r.GetAlive() // el jugador debe estar en la posicion 0 de la lista
					play_int, err32 := strconv.Atoi(play)
					if err32 != nil {
						log.Fatalf("fallo 32: %v", err32)
					}
					list_of_players[0].score = list_of_players[0].score + int32(play_int)
				
			} else {
				fmt.Println("el jugador está muerto")
			}
			// sección bots

			//Este pedazo de código es para las jugadas de los bots
			for i := 1; i < 16; i++ {

				if list_of_players[i].alive == true {
						fmt.Println(strconv.Itoa(i))
						botPlayer := list_of_players[i].id
						conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

						if err != nil {
							panic("cannot connect with server " + err.Error())
						}

						servicePlayer := pb.NewSquidGameServiceClient(conn)

						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						//jugada aleatoria
						rand.Seed(time.Now().UnixNano())
						playsend := rand.Int63n(9) + 1
						playsend_str := strconv.Itoa(int(playsend))
						r, err := servicePlayer.SendPlays(ctx, &pb.SendRequest{Player: botPlayer, Play: playsend_str, Stage: actualStage, Round: list_of_players[i].round})
						if err != nil {
							log.Fatalf("fallo 1: %v", err)
						}
						actualStage = r.GetStage()
						list_of_players[i].round = r.GetRound()
						list_of_players[i].alive = r.GetAlive()
						list_of_players[i].score = list_of_players[i].score + int32(playsend)
					
				}
				fmt.Println(list_of_players[i])
			}

		} else if action == "check" {
			message := "solicito monto"

			conn, err := grpc.Dial("10.6.43.41:8080", grpc.WithInsecure())

			if err != nil {
				panic("cannot connect with server " + err.Error())
			}

			servicePlayer := pb.NewSquidGameServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := servicePlayer.AmountCheck(ctx, &pb.AmountRequest{Message: message})
			if err != nil {
				log.Fatalf("no se pudo solicitar el monto: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMonto())

		} else {
			fmt.Println("ingresaste mal el comando")
		}
	}
}