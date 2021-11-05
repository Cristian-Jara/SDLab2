package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
	
)

type server struct {
	pb.UnimplementedChatServiceServer
}
type PlayerData struct{
	player string
	paths []string
}

func AppendData(player string, path string){
	for _,info := range PlayersData {
		if info.player == player {
			info.paths = append(info.paths, path)
		}
	}
}

func ReadData(player string)([]string){
	for _,info := range PlayersData {
		if info.player == player {
			return info.paths
		}
	}
	return nil
}

func (s *server) SendPlays(ctx context.Context, in *pb.SendPlay) (*pb.SendResult, error) {
	//aqui implementar la escritura del archivo de texto
	var path = "DataNode/Jugadas/jugador_" + in.Player + "_ronda_" + in.Stage + ".txt"
	AppendData(in.Player,path)
	//Verifica que el archivo existe
	var _, err = os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return &pb.SendResult{ Alive: false },err
		}
		defer file.Close()
	}

	// añadir al texto
	b, errtxt := ioutil.ReadFile(path)

	if errtxt != nil {
		log.Fatal(errtxt)
	}

	b = append(b, []byte(strconv.Itoa(int(in.Plays))+" \n")...)
	errtxt = ioutil.WriteFile(path, b, 0644)

	if errtxt != nil {
		log.Fatal(errtxt)
	}

	fmt.Println("Jugada del Jugador "+ in.Player+ " en la etapa "+ in.Stage +" recibida")
	return &pb.SendResult{ Alive: true },nil
}

func (s *server) GetPlayerInfo(ctx context.Context, in *pb.PlayerInfo) (*pb.PlayerInfo, error) {
	paths := ReadData(in.Message)
	if paths == nil{
		return &pb.PlayerInfo{Message: ""}, nil
	}
	message := ""
	for _, path := range paths {
		//Leer el archivo y chantar todo
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		message += string(b)
	}
	return &pb.PlayerInfo{Message: message}, nil
}
var PlayersData []PlayerData

func main() {
	var path = "DataNode/Jugadas"
	var _, err = os.Stat(path)
	if os.IsExist(err){
		os.RemoveAll(path) //Si existe se borra para no guardar datos de juegos anteriores
	}
	if _,err := os.Stat(path); os.IsNotExist(err){
		err = os.Mkdir(path, 0755)
		if err != nil{
			log.Fatalf("Failed to create the directory: %v",err)
		}
	}
	listner, err := net.Listen("tcp", ":50058")
	i := 1
	for i<=16{
		PlayersData = append(PlayersData, PlayerData{strconv.Itoa(i),nil})
	}
	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	servDN := grpc.NewServer()
	pb.RegisterChatServiceServer(servDN, &server{})

	//esto es lo que estaba al final, no sé donde ponerlo
	if err = servDN.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
