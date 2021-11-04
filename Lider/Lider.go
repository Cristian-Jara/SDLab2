package main

import(
	"log"
	"net"
	"fmt"
	"context"
	"math/rand"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
}
/*
func (s *Server) SendMessage(ctx context.Context, message *pb.Message)(* pb.Message, error){
	if GameStart != true {
		PlayersCount +=1
		if PlayersCount == 16 {
			GameStart = true
		}
		log.Printf("Received message body from client: %s",message.Body)
		log.Printf("Se recibio el jugador n° %v",PlayersCount)
		return &pb.Message{Body: fmt.Sprint("Hello from the server Jugador ",PlayersCount)}, nil
	} else {
		if message.Body == "Play" {
			return &pb.Message{Body: "Lo siento los jugadores máximos se han alcanzado"}, nil
		}
		log.Printf("Se esta jugando, %s", ctx)
		return &pb.Message{Body: "Se esta jugando"}, nil
	}
}
*/

func (s *Server)  JoinToGame(ctx context.Context, message *pb.JoinRequest) (*pb.JoinReply, error){
	if GameStart != true && message.Request == "Play"{
		PlayersCount +=1
		if PlayersCount == PlayerLimit {
			GameStart = true
			Start = true
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
		return &pb.GameStarted{Body: fmt.Sprint("Si"), Type: TypeStart},nil
	}else {
		return &pb.GameStarted{Body: "No"},nil
	}
} 
func (s *Server)  PlayTheGame(ctx context.Context, message *pb.Play)(*pb.Result,error){
	//if stage 1
	log.Printf("Received message body from client: %v",message.Plays)
	if int(message.Plays) < (rand.Intn(5) + 6) {
		return &pb.Result{RoundResult:"Sobreviviste", StageStatus:"Aún no terminada"}, nil
	}else{
		return &pb.Result{RoundResult:"Moriste", StageStatus:"Aún no terminada"}, nil
	}
}
func (s *Server)  GetMoneyAmount(ctx context.Context, message *pb.MoneyAmount)(*pb.MoneyAmount,error){
	return &pb.MoneyAmount{Money:"0"},nil
}
const (
	puerto = ":50052"
	PlayerLimit = 1
)
var ( 
	PlayersCount = 0
	RoundCount = 0
	GameStart = false
	Start = false
	TypeStart = "Stage 1"
)

func main(){
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
}