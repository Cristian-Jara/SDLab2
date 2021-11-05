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
}


func (s *Server) JoinToGame(ctx context.Context, message *pb.JoinRequest) (*pb.JoinReply, error){
	if GameStart != true && message.Request == "Play"{
		PlayersCount +=1
		if PlayersCount == PlayerLimit {
			GameStart = true
		}
		/////////////////////////////////////////////////////////////////////////
		playerlist = append(playerlist, PlayerNode{strconv.Itoa(PlayersCount), 0, true})
		/////////////////////////////////////////////////////////////////////////

		log.Printf("Received message body from client: %s",message.Request)
		log.Printf("Se recibio el jugador n° %v",PlayersCount)
		return &pb.JoinReply{Reply: fmt.Sprint("Bienvenido al Squid Game, eres el Jugador ",PlayersCount),Player: fmt.Sprint(PlayersCount)}, nil
	} else {
		return &pb.JoinReply{Reply: "Lo siento los jugadores máximos se han alcanzado", Player: "-1"}, nil
	}
} 
func (s *Server)  StageOrRoundStarted(ctx context.Context, message *pb.GameStarted) (*pb.GameStarted, error){
	if Start == true {
		leftPlayers-=1
		if leftPlayers == 0 {
			Start = false
		}
		return &pb.GameStarted{Body: "Si", Type: Stage},nil
	}else {
		return &pb.GameStarted{Body: "No"},nil
	}
} 

func JugadorMuerto(Body string){
	conn, err := amqp.Dial("amqp://admin:test@10.6.40.227:50069/")
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
	Sobrevive := true
	leftPlayers-=1
	if leftPlayers == 0 {
		Start = false
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
			return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil
		}else{
			Sobrevive = false
			//JugadorMuerto("Jugador_"+message.Player+" Ronda_"+Stage)
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
			time.Sleep(10*time.Millisecond)
		} 
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
				for _,aux := range team1{
					aux.live = false
				}
				teamf = team2
			}else{
				for _, aux := range team2{
					aux.live = false
				}
				teamf = team1
			}
		} else if jugada%2==score_team1%2 && jugada%2!=score_team2%2 { // TEAM 1 WINS
			for _,aux := range team2{
				aux.live = false
			}
			teamf = team1
		} else if jugada%2==score_team1%2 && jugada%2==score_team2%2 { // TEAM 2 WINS
			for _,aux := range team1{
				aux.live = false
			}
			teamf = team2
		} // IF EVERYBODY WINS NOBODY DIES
		for idx,_ := range team1{
			if team1[idx].id == message.Player{
				return &pb.SendResult{Stage: Stage, Alive: team1[idx].live, Round: int32(RoundCount), Started: true}, nil
			}else if team2[idx].id == message.Player{
				return &pb.SendResult{Stage: Stage, Alive: team1[idx].live, Round: int32(RoundCount), Started: true}, nil
			}
		}
		return nil, nil
	}else if Stage =="3"{
		log.Printf("Jugando etapa 3")
		for idx,_ := range teamgg{
			if teamgg[idx].id == message.Player{
				teamgg[idx].score = int(message.Plays)
			}
		}

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
			i+=2
		}
		for idx,_ := range teamgg{
			if teamgg[idx].live{
				fmt.Println("El jugador: " + teamgg[idx].id + " ha ganado el Squid Game")
			}
			return &pb.SendResult{Stage: Stage, Alive: teamgg[idx].live, Round: int32(RoundCount), Started: true}, nil
		}
		return nil, nil
		//return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil

	}else{
		log.Printf("Termino el juego")
		return &pb.SendResult{Stage: Stage, Alive: true, Round: int32(RoundCount), Started: true}, nil

	}
}
func (s *Server)  GetMoneyAmount(ctx context.Context, message *pb.MoneyAmount)(*pb.MoneyAmount,error){
	conn, err := grpc.Dial(":50065", grpc.WithInsecure()) //10.6.40.225
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
	puerto = ":50052"
	PlayerLimit = 2 // Cantidad de jugadores aceptados
	NameNodeAddress = ":50055" //"10.6.40.227:50055"
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
		/*
		for input != "start" || leftPlayers != 0{
			fmt.Println("Ingresa 'start' para comenzar la etapa 1: ")
			fmt.Scanln(&input)
			if input == "start" { 
				Start = true
				fmt.Println("Esperando que todos los jugadores esten listos ....")
				for leftPlayers!=0{
					continue
				}
			}
		}*/
		interfaz(Stage)
		fmt.Println("Ha comenzado la etapa: " + Stage)
		rand.Seed(time.Now().UnixNano())
		for RoundCount < 5 {
			time.Sleep(1*time.Second)
			jugada = rand.Intn(5) + 6
			fmt.Println("Jugada de Lider: " + strconv.Itoa(jugada))
			leftPlayers = PlayersCount // Contar que los jugadores ingresen a la ronda
			Start = true
			for leftPlayers != 0 {
				time.Sleep(10*time.Millisecond)
			}
			leftPlayers = PlayersCount //Contar que los jugadores jueguen
			for leftPlayers != 0 {
				time.Sleep(10*time.Millisecond)
			}
			RoundCount+=1
		}
			
		for j := 0; j < PlayerLimit; j++ {
			if (playerlist[j].score < 21) && (playerlist[j].live == true) {
				playerlist[j].live = false
				puntaje := strconv.Itoa(int(playerlist[j].score))
				fmt.Println("El jugador: " + playerlist[j].id + " fue eliminado con puntaje : " + puntaje)
				statement := "Jugador_" + playerlist[j].id + " Ronda_" + Stage
				log.Printf(" Ha muerto: %s ", statement)
			}
		}

		//// Jugadores ya eliminados, se anuncian los sobrevivientes  /////////

		winners = 0
		for i := 0; i < PlayerLimit; i++ {
			playerlist[i].score = 0
			if playerlist[i].live == true {
				winners += 1
				fmt.Println("Jugador: " + playerlist[i].id + " pasa el nivel 1")
			}
		}
		Start = false
		Stage = "2"
		interfaz(Stage)
		fmt.Println("Ha comenzado la etapa: " + Stage)
		for winners%2 == 1{
			//Borrar un jugador al azar porque sobra 1//
			jugada = rand.Intn(15)
			if playerlist[jugada].live == true {
				playerlist[jugada].live = false
				winners -= 1
				fmt.Println("Jugador: " + playerlist[jugada].id + " es eliminado al azar")
				statement := "Jugador_" + playerlist[jugada].id + " Ronda_" + Stage
				JugadorMuerto("Jugador_"+playerlist[jugada].id+" Ronda_"+Stage)
				log.Printf(" Ha muerto: %s ", statement)
			}
		}
		
		// Empieza la etapa 2 separando los ganadores de la ronda 1 en 2 grupos //
		swap := 0
		for i := 0; i < 16; i++ {
			if playerlist[i].live == true {
				if swap == 0 {
					team1 = append(team1, PlayerNode{playerlist[i].id, 0, true})
					fmt.Println("Se agrega a team 1: " + playerlist[i].id)
					swap = 1
				} else {
					team2 = append(team2, PlayerNode{playerlist[i].id, 0, true})
					fmt.Println("Se agrega a team 2: " + playerlist[i].id)
					swap = 0
				}

			}
		}

		// Hasta acá se tienen los jugadores restantes en 2 equipos, se inicia etapa 2 //

		rand.Seed(time.Now().UnixNano())
		jugada = rand.Intn(3)
		jugada = jugada + 1
		fmt.Println("Jugada del lider: " + strconv.Itoa(jugada))
		for leftPlayers != 0 {
				time.Sleep(10*time.Millisecond)
		}
		time.Sleep(1*time.Second)
		//fmt.Println("Ingresa start tras terminar los turnos de los jugadores: ")
		//fmt.Scanln(&start)
		/*
		t1score := 0
		t2score := 0
		t1win := false
		t2win := false
		for i := 0; i < len(team1); i++ {
			t1score += team1[i].score
		}
		for i := 0; i < len(team2); i++ {
			t2score += team2[i].score
		}
		fmt.Println("Score Team 1: " + strconv.Itoa(t1score))
		fmt.Println("Score Team 2: " + strconv.Itoa(t2score))

		if t1score%2 == jugada%2 {
			t1win = true
		}
		if t2score%2 == jugada%2 {
			t2win = true
		}

		if t1win == true && t2win == true {
			fmt.Println("¡Ambos equipos avanzan a la ronda final!")
			winners = len(team1) + len(team2)
			teamf = append(team1, team2)
		}
		else if t1win == true && t2win == false {
			fmt.Println("Team 1 avanza a la ronda final")
			for i := 0; i < len(team2); i++ {
				team2[i].live = false
			}
			winners = len(team1)
			teamf = team1
		}
		else if t1win == false && t2win == true {
			fmt.Println("Team 2 avanza a la ronda final")
			for i := 0; i < len(team1); i++ {
				team1[i].live = false
			}
			winners = len(team2)
			teamf= team2
		}
		else {
			fmt.Println("Ambos equipos pierden, se elimina uno al azar para continuar")
			rand.Seed(time.Now().UnixNano())
			jugada = rand.Intn(1))
			if jugada == 0 {
				fmt.Println("Team 1 avanza a la ronda final")
				for i := 0; i < len(team2); i++ {
					team2[i].live = false
				}
				winners = len(team1)
				teamf = team1
			} else {
				fmt.Println("Team 2 avanza a la ronda final")
				for i := 0; i < len(team1); i++ {
					team1[i].live = false
				}
				winners = len(team2)
				teamf = team2
			}
		}*/
		
		//////////// Aca termina la ronda 2, se anuncian los jugadores vivos ///////////////////

		fmt.Println("Los jugadores que pasan a la siguiente ronda son en total: " + strconv.Itoa(len(teamf)))
		for i := 0; i < len(teamf); i++ {
			teamf[i].score = 0
			if teamf[i].live == true {
				fmt.Println("El jugador: " + teamf[i].id + " pasa a la siguiente etapa")
			}
		}

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
		Stage = "3"
		Start = false

		////////////////// Acá empieza la etapa 3 y final///////////////////////////////////
		/*
		for i := 0; i < len(teamf); i++ {
			if teamf[i].live == true {
				teamgg = append(teamgg, PlayerNode{teamf[i].id, 0, true})
			}
			Start = false
			interfaz(Stage)
			jugada = rand.Intn(9))
			jugada = jugada+ 1
			fmt.Println("Jugada de lider: " + strconv.Itoa(jugada))

			////// Loop de jugadas para cada player ///////
			for i := 0; i < len(teamgg); i++ {
				teamgg[i].score = rand.Intn(9)) + 1
				i++
			}

			////// Se evaluan el valor absoluto de los puntajes //////
			for i := 0; i < len(teamgg); i++ {
				teamgg[i].score = math.Abs(int64(teamgg[i].score) - int64(jugada))
				i++
			}

			// Se itera sobre la lista final eligiendo ganadores en competencias de a pares //
			swap = 0
			score1 := 0
			score2 := 2
			for i := 0; i < len(teamgg); i++ {
				score1 = teamgg[i].score
				score2 = teamgg[i+1].score
				if score1 > score2{
					teamgg[i+1].live = false
					statement := "Jugador_" + teamgg[i+1].id + " Ronda_" + Stage
					log.Printf(" Ha muerto: %d ", statement)
				} else if score1 < score2{
					teamgg[i].live = false
					statement := "Jugador_" + teamgg[i].id + " Ronda_" + Stage
					log.Printf(" Ha muerto: %d ", statement)
				}
				i+=2
			}

			// Teniendo la lista final de vencedores, se anuncia quienes ganaron ///
			for i := 0; i < len(teamgg); i++ {
				if teamgg[i].live == true{
					fmt.Println("El jugador: " + teamgg[i].id + " ha ganado el Squid Game")
				}
				i++
			}


		}*/
		

				




	}
	
	for {
		fmt.Scanln(&input)
		if input == "end" {
			break
		} else if input == "play"{
			Start = true
			leftPlayers = PlayersCount 
		} else if input == "stop" {
			Start = false
		}
	}
}