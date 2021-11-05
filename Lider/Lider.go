package main

import(
	"log"
	"net"
	"fmt"
	"context"
	"math/rand"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
}

func (s *Server) JoinToGame(ctx context.Context, message *pb.JoinRequest) (*pb.JoinReply, error){
	if GameStart != true && message.Request == "Play"{
		PlayersCount +=1
		if PlayersCount == PlayerLimit {
			GameStart = true
		}
		log.Printf("Received message body from client: %s",message.Request)
		log.Printf("Se recibio el jugador n° %v",PlayersCount)
		return &pb.JoinReply{Reply: fmt.Sprint("Bienvenido al Squid Game, eres el Jugador ",PlayersCount),Player: fmt.Sprint(PlayersCount)}, nil
	} else {
		return &pb.JoinReply{Reply: "Lo siento los jugadores máximos se han alcanzado", Player: "-1"}, nil
	}
} 
func (s *Server)  StageOrRoundStarted(ctx context.Context, message *pb.GameStarted) (*pb.GameStarted, error){
	if Start == true {
		return &pb.GameStarted{Body: fmt.Sprint("Si"), Type: Stage},nil
	}else {
		return &pb.GameStarted{Body: "No"},nil
	}
} 

func JugadorMuerto(Body string){
	conn, err := amqp.Dial("amqp://admin:test@localhost:50069/")
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
}

func (s *Server)  SendPlays(ctx context.Context, message *pb.SendPlay)(*pb.SendResult,error){
	//if stage 1
	Sobrevive := true
	log.Printf("Received message body from client: %v",message.Plays)
	conn, err := grpc.Dial(NameNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	NameNodeServer := pb.NewChatServiceClient(conn)
	NameNodeServer.SendPlays(context.Background(), &pb.SendPlay{Player: message.Player, Plays: message.Plays,  Stage: message.Stage, Round: message.Round, Score: message.Score})
	if Stage == "1"{
		if int(message.Plays) < (rand.Intn(5) + 6) {
			return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil
		}else{
			Sobrevive = false
			JugadorMuerto("Jugador_"+message.Player+" Ronda_"+Stage)
			return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil
		}
	}else if Stage =="2"{
		log.Printf("Jugando etapa 2")
		return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil

	}else if Stage =="3"{
		log.Printf("Jugando etapa 3")
		return &pb.SendResult{Stage: Stage, Alive: Sobrevive, Round: int32(RoundCount), Started: true}, nil

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
	PlayerLimit = 2
	NameNodeAddress = ":50055" //"10.6.40.227:50055"
)
var ( 
	PlayersCount = 0
	RoundCount = 0
	GameStart = false
	Start = false
	Stage = "1"
	input string
)

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


	fmt.Println("Esta es la interfaz para el lider")
	for {
		fmt.Scanln(&input)
		if input == "end" {
			break
		} else if input == "play"{
			Start = true
		} else if input == "stop" {
			Start = false
		}
	}
}