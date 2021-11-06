package main

import(
	"log"
	"net"
	"fmt"
	"context"
	"math/rand"
	"strconv"
	"time"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
}	

type PlayerNode struct {
	id    string
	score int
	live bool
	played bool
}


func (s *Server) JoinToGame(ctx context.Context, message *pb.JoinRequest) (*pb.JoinReply, error){
	if GameStart != true && message.Request == "Play"{
		PlayersCount +=1
		if PlayersCount == PlayerLimit {
			GameStart = true
		}
		/////////////////////////////////////////////////////////////////////////
		playerlist = append(playerlist, PlayerNode{strconv.Itoa(PlayersCount), 0, true, false})
		/////////////////////////////////////////////////////////////////////////

		log.Printf("Received message body from client: %s",message.Request)
		log.Printf("Se recibio el jugador n° %v",PlayersCount)
		return &pb.JoinReply{Reply: fmt.Sprint("Bienvenido al Squid Game, eres el Jugador ",PlayersCount),Player: fmt.Sprint(PlayersCount)}, nil
	} else {
		return &pb.JoinReply{Reply: "Lo siento los jugadores máximos se han alcanzado", Player: "-1"}, nil
	}
} 
func (s *Server)  StageOrRoundStarted(ctx context.Context, message *pb.GameStarted) (*pb.GameStarted, error){
	id,_ := strconv.Atoi(message.Body)
	if playerlist[id -1].live == true{
		if GameFinished == false{
			if Start == true && playerlist[id -1].played == false {
				playerlist[id - 1].played = true
				leftPlayers-=1
				if leftPlayers == 0 {
					Start = false
					for idx, _ := range playerlist {
						playerlist[idx].played = false
					}
				}
				return &pb.GameStarted{Body: "Si", Type: Stage},nil
			}else {
				return &pb.GameStarted{Body: "No"},nil
			}
		}else{
			PlayerKnows = true
			return &pb.GameStarted{Body: "Win"},nil
		}
	}else{
		return &pb.GameStarted{Body: "Killed"},nil
	}
} 

func JugadorMuerto(Body string){
	conn, err := amqp.Dial("amqp://guest:guest@10.6.40.227:50069/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(Body),
		})
	if err != nil {
		log.Fatalf("Failed to publish the message: %v", err)
	}
	PlayersCount -=1
}

func (s *Server)  SendPlays(ctx context.Context, message *pb.SendPlay)(*pb.SendResult,error){
	//if stage 1
	id,_ := strconv.Atoi(message.Player)
	playerlist[id - 1].played = true
	Sobrevive := true
	leftPlayers-=1
	if leftPlayers == 0 {
		Start = false
		for idx, _ := range playerlist {
			playerlist[idx].played = false
		}
	}
	log.Printf("Received message body from client: %v",message.Plays)
	conn, err := grpc.Dial(NameNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	NameNodeServer := pb.NewChatServiceClient(conn)
	NameNodeServer.SendPlays(context.Background(), &pb.SendPlay{Player: message.Player, Plays: message.Plays,  Stage: message.Stage, Round: message.Round, Score: message.Score})
	if Stage == "1"{
		if int(message.Plays) < jugada {
			playerlist[id - 1].score += int(message.Plays)
			return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil
		}else{
			Sobrevive = false
			PlayersCount-=1
			playerlist[id - 1].score += int(message.Plays)
			playerlist[id - 1].live = false
			JugadorMuerto("Jugador_"+message.Player+" Ronda_"+Stage)
			return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil
		}
	}else if Stage =="2"{
		for idx,_ := range team1{
			if team1[idx].id == message.Player{
				team1[idx].score = int(message.Plays)
			}else if team2[idx].id == message.Player{
				team2[idx].score = int(message.Plays)
			}
		}
		for leftPlayers != 0 {
			time.Sleep(20*time.Millisecond)
		}
		time.Sleep(100*time.Millisecond)
		score_team1 := 0
		for _,aux := range team1{
			score_team1 += aux.score
		}
		score_team2 := 0
		for _,aux:= range team2{
			score_team2 += aux.score
		}
		if jugada%2!=score_team1%2 && jugada%2!=score_team2%2 { // NOBODY WINS
			die:=rand.Intn(2)
			if die==0{
				for idx,_ := range team1{
					team1[idx].live = false
				}
				teamf = team2
			}else{
				for idx,_ := range team2{
					team2[idx].live = false
				}
				teamf = team1
			}
		} else if jugada%2==score_team1%2 && jugada%2!=score_team2%2 { // TEAM 1 WINS
			for idx,_ := range team2{
				team2[idx].live = false
			}
			teamf = team1
		} else if jugada%2!=score_team1%2 && jugada%2==score_team2%2 { // TEAM 2 WINS
			for idx,_ := range team1{
				team1[idx].live = false
			}
			teamf = team2
		} else{ // IF EVERYBODY WINS NOBODY DIES
			return &pb.SendResult{Stage: Stage, Alive: true, Round: int32(RoundCount), Started: true}, nil
		}
		Sobrevive := false
		for _,value := range teamf{
			if value.id == message.Player{
				Sobrevive = true
			}
		}
		if Sobrevive != true {
			playerlist[id - 1].live = false
			JugadorMuerto("Jugador_"+message.Player+" Ronda_"+Stage)
		}
		return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil
	}else if Stage =="3"{
		for idx,_ := range teamgg{
			if teamgg[idx].id == message.Player{
				teamgg[idx].score = int(message.Plays)
			}
		}

		for leftPlayers != 0 {
			time.Sleep(20*time.Millisecond)
		}
		time.Sleep(100*time.Millisecond)

		score_j1 := 0
		score_j2 := 0
		for i := 0; i < len(teamgg); i++ {
			score_j1 = teamgg[i].score-jugada
			score_j2 = teamgg[i+1].score-jugada
			if score_j1<0{
				score_j1 = score_j1*-1
			}else if score_j2<0{
				score_j2 = score_j2*-1
			}

			if score_j1 > score_j2{
				teamgg[i+1].live = false
			} else if score_j1 < score_j2{
				teamgg[i].live = false
			}
			i+=1
		}

		Sobrevive := false
		for _,value := range teamgg{
			if value.id == message.Player{
				Sobrevive = true
			}
		}
		if Sobrevive != true{
			playerlist[id - 1].live = false
			JugadorMuerto("Jugador_"+message.Player+" Ronda_"+Stage)
		}

		return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil

		//return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil

	}else{
		log.Printf("Termino el juego")
		return &pb.SendResult{Stage: Stage, Alive: true, Round: int32(RoundCount), Started: true}, nil

	}
}
func (s *Server)  GetMoneyAmount(ctx context.Context, message *pb.MoneyAmount)(*pb.MoneyAmount,error){
	conn, err := grpc.Dial("10.6.40.225:50065", grpc.WithInsecure()) //10.6.40.225
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	sc := pb.NewChatServiceClient(conn)
	r, err := sc.GetMoneyAmount(context.Background(), &pb.MoneyAmount{Money: ""})
	if err != nil {
		log.Fatalf("No se pudo obtener el monto: %v", err)
	}
	return &pb.MoneyAmount{Money: r.Money}, nil
}
const (
	puerto = ":50052" // "xd"
	PlayerLimit = 2 // Cantidad de jugadores aceptados
	NameNodeAddress = "10.6.40.227:50055"
)
var ( 
	leftPlayers = 2 // Jugadores por entrar a jugar
	next string
	playerlist []PlayerNode	
	jugada = 0
	PlayersCount = 0 // Cantidad de jugadores actuales
	RoundCount = 1
	GameStart = false
	Start = false
	Stage = "1"
	input string
	winners = 0
	GameFinished = false
	PlayerKnows = false
)

var team1 []PlayerNode
var team2 []PlayerNode
var teamf []PlayerNode
var teamgg []PlayerNode

func interfaz(Stage string){
	for{
		fmt.Println("Ingresa 'start' para comenzar la etapa "+Stage+" o\nIngresa 'view' para ver las jugadas de un cierto jugador")
		fmt.Scanln(&input)
		for input != "start" && input != "view"{
			fmt.Println("Opción no válida, ingrese nuevamente")
			fmt.Scanln(&input)
		}
		if input == "start"{
			break
		}
		if input == "view" {
			fmt.Println("Ingrese el número del jugador a buscar las jugadas")
			var PlayerToSearch int
			fmt.Scanln(&PlayerToSearch)
			for PlayerToSearch > PlayerLimit || PlayerToSearch < 1 {
				fmt.Println("Ingrese un número de jugador válido")
				fmt.Scanln(&PlayerToSearch)
			}
			conn, err := grpc.Dial(NameNodeAddress, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			NameNodeServer := pb.NewChatServiceClient(conn)
			r, err := NameNodeServer.GetPlayerInfo(context.Background(), &pb.PlayerInfo{Message: strconv.Itoa(PlayerToSearch)})
			if err != nil {
				log.Fatalf("Error al obtener los datos del jugador")
			}
			fmt.Println("Las jugadas realizadas hasta ahora por el jugador son las siguientes:\n" + r.Message)
		}
	}
}



func main(){
	go func () {
		lis, err := net.Listen("tcp", puerto)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterChatServiceServer(s, &Server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}() // Este es el servidor con el cual se hacen las comunicaciones

	time.Sleep(1*time.Second)
	for PlayersCount != PlayerLimit{
		fmt.Println("Ingresa 'start' para iniciar el juego: ")
		fmt.Scanln(&input)
		// Empieza la fase 1 de SG
		if PlayersCount != PlayerLimit{
			fmt.Println("Se requieren mas jugadores para iniciar Squid Game")
		}
	}

	if PlayersCount == PlayerLimit{   // acá empieza sg
		interfaz(Stage)
		fmt.Println("Ha comenzado la etapa: " + Stage)
		rand.Seed(time.Now().UnixNano())
		Start = true
		fmt.Println("Esperando que los jugadores ingresen a la etapa ...")
		for leftPlayers != 0 { //Contar que los jugadores ingresarán a la etapa
			time.Sleep(20*time.Millisecond)
		}
		for RoundCount < 5 {
			time.Sleep(1*time.Second)
			jugada = rand.Intn(5) + 6
			fmt.Println("Jugada de Lider: " + strconv.Itoa(jugada))
			leftPlayers = PlayersCount // Contar que los jugadores ingresen a la ronda
			Start = true
			fmt.Println("Esperando que todos los jugadores ingresen a la ronda ...")
			for leftPlayers != 0 {
				time.Sleep(20*time.Millisecond)
			}
			leftPlayers = PlayersCount //Contar que los jugadores jueguen
			Start = true
			fmt.Println("Esperando que todos los jugadores ingresen su jugada ...")
			for leftPlayers != 0 {
				time.Sleep(20*time.Millisecond)
			}
			RoundCount+=1
		}
			
		for j := 0; j < PlayerLimit; j++ {
			if (playerlist[j].score < 21) && (playerlist[j].live == true) {
				playerlist[j].live = false
				PlayersCount -=1
				puntaje := strconv.Itoa(int(playerlist[j].score))
				fmt.Println("El jugador: " + playerlist[j].id + " fue eliminado con puntaje : " + puntaje)
				statement := "Jugador_" + playerlist[j].id + " Ronda_" + Stage
				log.Printf(" Ha muerto: %s ", statement)
				JugadorMuerto("Jugador_"+playerlist[j].id +" Ronda_"+Stage)
			}
		}
		Start = false

		//// Jugadores ya eliminados, se anuncian los sobrevivientes  /////////

		winners = 0
		for i := 0; i < PlayerLimit; i++ {
			playerlist[i].score = 0
			if playerlist[i].live == true {
				winners += 1
				fmt.Println("Jugador: " + playerlist[i].id + " pasa el nivel 1")
			}
		}
		if PlayersCount == 1{
			for _,value := range playerlist{
				if value.live == true{
					fmt.Println("El jugador: " + value.id + " ha ganado el Squid Game")
					GameFinished = true
					for PlayerKnows != true{
						time.Sleep(1*time.Second)
					}
					return
				}
			}
		}
		Start = false
		Stage = "2"
		interfaz(Stage)
		fmt.Println("Ha comenzado la etapa: " + Stage)
		for winners%2 == 1{
			//Borrar un jugador al azar porque sobra 1//
			jugada = rand.Intn(PlayerLimit - 1)
			if playerlist[jugada].live == true {
				playerlist[jugada].live = false
				PlayersCount-=1
				winners -= 1
				fmt.Println("Jugador: " + playerlist[jugada].id + " es eliminado al azar")
				statement := "Jugador_" + playerlist[jugada].id + " Ronda_" + Stage
				JugadorMuerto("Jugador_"+playerlist[jugada].id+" Ronda_"+Stage)
				log.Printf(" Ha muerto: %s ", statement)
			}
		}

		
		// Empieza la etapa 2 separando los ganadores de la ronda 1 en 2 grupos //
		swap := 0
		for i := 0; i < PlayerLimit; i++ {
			if playerlist[i].live == true {
				if swap == 0 {
					team1 = append(team1, PlayerNode{playerlist[i].id, 0, true,false})
					fmt.Println("Se agrega a team 1: " + playerlist[i].id)
					swap = 1
				} else {
					team2 = append(team2, PlayerNode{playerlist[i].id, 0, true,false})
					fmt.Println("Se agrega a team 2: " + playerlist[i].id)
					swap = 0
				}

			}
		}
		leftPlayers = PlayersCount
		Start = true
		fmt.Println("Esperando que los jugadores ingresen a la etapa ...")
		for leftPlayers != 0 { //Contar que los jugadores ingresarán a la etapa
			time.Sleep(20*time.Millisecond)
		}

		// Hasta acá se tienen los jugadores restantes en 2 equipos, se inicia etapa 2 //
		jugada = rand.Intn(4) + 1
		fmt.Println("Jugada del lider: " + strconv.Itoa(jugada))
		leftPlayers = PlayersCount //Contar que los jugadores jueguen
		Start = true
		fmt.Println("Esperando que todos los jugadores ingresen su jugada ...")
		for leftPlayers != 0 {
			time.Sleep(20*time.Millisecond)
		}
		time.Sleep(1*time.Second)
		//////////// Aca termina la ronda 2, se anuncian los jugadores vivos ///////////////////

		fmt.Println("Los jugadores que pasan a la siguiente ronda son en total: " + strconv.Itoa(len(teamf)))
		j:= 0 
		for _,value := range playerlist{
			if value.live == true {
				fmt.Println("El jugador: " + value.id + " pasa a la siguiente etapa")
				j+=1
			}
		}
		PlayersCount = j
		////////////////// Se elimina el jugador que sobra, si hay uno //////////////////////
		/*
		for winners%2 == 1 {
			rand.Seed(time.Now().UnixNano())
			jugada = rand.Intn(15)
			if teamf[jugada].live == true {
				teamf[jugada].live = false
				winners -= 1
				fmt.Println("El jugador: " + teamf[jugada].id + " es eliminado aleatoriamente")
				statement := "Jugador_" + teamf[jugada].id + " Ronda_" + Stage
				log.Printf(" Ha muerto: %s ", statement)
			}
		}*/

		if PlayersCount == 1{
			for _,value := range playerlist{
				if value.live == true{
					fmt.Println("El jugador: " + value.id + " ha ganado el Squid Game")
					GameFinished = true
					for PlayerKnows != true{
						time.Sleep(1*time.Second)
					}
					return
				}
			}
		}
		Stage = "3"   
		            // -----------------------------------------------------------------------> LA MARIA AÑADIO DESDE ACA
		Start = false
		interfaz(Stage)
		fmt.Println("Ha comenzado la etapa: " + Stage)
		i:= 1
		v:= 1
		for _,value := range playerlist{
			if value.live == true{
				teamgg = append(teamgg,PlayerNode{value.id,0,true,false})
				fmt.Println("Jugador "+value.id+" en equipo "+strconv.Itoa(i))
				if v%2==0{
					i+=1
				}
				v+=1
			}
		}
		leftPlayers=PlayersCount
		Start = true
		fmt.Println("Esperando que los jugadores ingresen a la etapa ...")
		for leftPlayers != 0 { //Contar que los jugadores ingresarán a la etapa
			time.Sleep(20*time.Millisecond)
		}

		jugada = rand.Intn(10) + 1
		fmt.Println("Jugada del lider: " + strconv.Itoa(jugada))
		leftPlayers = PlayersCount //Contar que los jugadores jueguen
		Start = true
		fmt.Println("Esperando que todos los jugadores ingresen su jugada ...")
		for leftPlayers != 0 {
			time.Sleep(20*time.Millisecond)
		}
		time.Sleep(3*time.Second)
		for idx,_ := range playerlist{
			if playerlist[idx].live == true{
				fmt.Println("El jugador: " + playerlist[idx].id + " ha ganado el Squid Game")
			}
		}
		time.Sleep(10*time.Second)
		

				




	}

}