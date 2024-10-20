package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	//connectando ao rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatal("Error connecting to RabbitMQ", err)
	}

	defer conn.Close()

	// abrindo um cannal de comunicacao
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error ao abrir o canal", err)
	}

	defer ch.Close()

	// declarando o exchange do tipo topic
	err = ch.ExchangeDeclare(
		"meu_ex", // Nome do exchange
		"direct", // Tipo topic
		true,     // Duravel
		false,    //Auto-declare
		false,    //Não e uma exchange interna
		false,    //No-wait
		nil,      //Sem argumento
	)

	if err != nil {
		log.Fatal("Error ao criar o exchange do tipo topic", err)
	}

	message := os.Args[1]

	// definindo a mensagem
	routingKey := "" // Valor padrão
	if len(os.Args) > 2 {
		routingKey = os.Args[2] // usar routing key se existir
	}

	// Publicando a messagem na exchange
	err = ch.Publish(
		"meu_ex",   //Nome do exchange
		routingKey, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	//failOnError(err, "Failed to declare a queue")

	if err != nil {
		log.Fatal("Error ao publicar")
	}

	fmt.Printf("Mensagem enviada : %s com a rpouting key %s\n", message, routingKey)
}
