package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	//connectando ao rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		failOnError("Error connecting to RabbitMQ", err)
	}
	defer conn.Close()
	// abrindo um cannal de comunicacao
	ch, err := conn.Channel()
	if err != nil {
		failOnError("Error ao abrir o canal", err)
	}

	defer ch.Close()
	// declarando a fila da qual o consumer vai receber
	q, err := ch.QueueDeclare(
		"log", // Nome da fila
		true,  // Duravel
		false, //Delete quando não estiver em uso
		false, //Exclusiva
		false, //No-wait
		nil,   //Sem argumento
	)

	if err != nil {
		failOnError("Error ao  declarar a fila", err)
	}

	// consumindo mensagem da fila
	msgs, err := ch.Consume(
		q.Name,       //Nome da fila
		"pagamentos", //Consumer
		true,         //Auto-ack CUIDADO
		false,        //Exclusiva
		false,        //No-local
		false,        //No-wait
		nil,          //Sem argumento
	)

	if err != nil {
		failOnError("Error ao  registar a fila", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf("Mensagem Recebida %s", msg.Body)
		}
	}()

	log.Println("Aguardando mensagens serem lidas ...")
	//Lendo o canal
	<-forever

}
