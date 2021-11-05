package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"log"
	"net"
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/Cristian-Jara/SDLab2.git/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) AmountCheck(ctx context.Context, in *pb.MoneyAmount) (*pb.MoneyAmount, error) {
	return &pb.MoneyAmount{Money: strconv.Itoa(DineroAcum)}, nil
}
func ErrorMessage(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var DineroAcum int = 0

const (
	LiderIP = "10.6.40.227"
	LocalIP = "localhost"
	Puerto = ":50069"
)

func main() {

	go func() {
		listner, err := net.Listen("tcp", ":50065")

		if err != nil {
			panic("cannot connect with server " + err.Error())
		}

		serv := grpc.NewServer()
		pb.RegisterChatServiceServer(serv, &server{})
		if err = serv.Serve(listner); err != nil {
			panic("cannot initialize the server" + err.Error())

		}
	}()

	conn, err := amqp.Dial(fmt.Sprint("amqp://admin:test@",LocalIP,Puerto,"/"))
	ErrorMessage(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	ErrorMessage(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	ErrorMessage(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	ErrorMessage(err, "Failed to register a consumer")
	forever := make(chan bool)
	path := "Pozo/registro_muertes.txt"
	_, err = os.Stat(path)
	if os.IsExist(err) {
		os.Remove(path) //Si existe se borra para no guardar datos de juegos anteriores
	}
	file, err := os.Create(path) //Se crea el archivo de registro
	if err != nil {
		log.Fatalf("failed to create the register: %v", err)
	}
	defer file.Close()
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			cadena := string(d.Body)
			DineroAcum = DineroAcum + 100000000
			b = append(b, []byte(cadena+" "+strconv.Itoa(DineroAcum)+" \n")...)
			err = ioutil.WriteFile(path, b, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
