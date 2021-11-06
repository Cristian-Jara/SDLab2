package main

import (
	"log"
	"fmt"
	"context"
	"io/ioutil"
	"math/rand"
	"net"
	"time"
	"os"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) SendPlays(ctx context.Context, in *pb.SendPlay) (*pb.SendResult, error) {
	//enviar la jugada a cualquiera de los 3.
	var address string
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(3)
	if id == 0 {
		address = "10.6.40.225" // IP1
	} else if id == 1 {
		address = "10.6.40.227" // IP2
	} else {
		address = "10.6.40.229" // IP3
	}
	fmt.Println("Nueva jugada recibida: \nJugador " + in.Player + " jugó:  %v",in.Plays)
	conn, err := grpc.Dial(address+":50058", grpc.WithInsecure())
	DataNodeService := pb.NewChatServiceClient(conn)
	r, err := DataNodeService.SendPlays(context.Background(), &pb.SendPlay{Player: in.Player, Plays: in.Plays, Stage: in.Stage})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Stage)
	// añadir al texto
	b, errtxt := ioutil.ReadFile("NameNode/registro.txt")

	if errtxt != nil {
		log.Fatal(errtxt)
	}

	b = append(b, []byte("Jugador_"+in.Player+" Ronda_"+in.Stage+" "+address+"\n")...)
	errtxt = ioutil.WriteFile("NameNode/registro.txt", b, 0644)

	if errtxt != nil {
		log.Fatal(errtxt)
	}

	return &pb.SendResult{}, nil
}

func (s *server) GetPlayerInfo(ctx context.Context, in *pb.PlayerInfo) (*pb.PlayerInfo, error) {
	list := [3]string{"10.6.40.225:50058","10.6.40.227:50058","10.6.40.229:50058"} //Ingresar las direcciones IP
	// "10.6.40.225:50058" // IP1
	// "10.6.40.227:50058" // IP2
	// "10.6.40.229:50058" // IP3
	message := ""
	for _, dir := range list {
		//Leer el archivo y chantar todo
		conn, err := grpc.Dial(dir, grpc.WithInsecure())
		DataNodeService := pb.NewChatServiceClient(conn)
		r, err := DataNodeService.GetPlayerInfo(context.Background(), &pb.PlayerInfo{Message: in.Message})
		if err != nil {
			log.Fatalf("Problema al obtener la información del jugador: %v", err)
		}
		message += r.Message
	}
	return &pb.PlayerInfo{Message: message}, nil
}

func main(){
	var path = "NameNode/registro.txt"
	var _, err = os.Stat(path)
	if os.IsExist(err) {
		os.Remove(path) //Si existe se borra para no guardar datos de juegos anteriores
	}
	file, err := os.Create(path) //Se crea el archivo de registro
	if err != nil {
		log.Fatalf("failed to create the register: %v", err)
	}
	defer file.Close()
	
	listner, err := net.Listen("tcp", ":50055")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	pb.RegisterChatServiceServer(serv, &server{})
	fmt.Println("Recibiendo jugadas desde el Lider en el puerto: 50055")
	if err = serv.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}