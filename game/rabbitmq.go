package game

import (
	"github.com/streadway/amqp"
	"log"
)

var (
	Logs				= false
	conn               	*amqp.Connection
	ch                 	*amqp.Channel
	rabbitMQURL        		= 	"amqp://tanks:ADf5KdBv8M91@35.184.207.230:5672/"
	//rabbitMQURL        		= 	"amqp://guest:guest@localhost:5672/"
	commandsExchange   		= 	"commands"
	onlineExchange			= 	"onlines"
	onlineQueue 		amqp.Queue
	receiverQueue      	amqp.Queue
)

func failOnError(err error, msg string, ok string) {
		if err != nil {
			log.Fatalf("%s: %s", msg, err)
		}
		if ok != "" && Logs{
			log.Printf(" [*] " + ok)
		}
}

func RabbitMQ()(err error){
	conn, err = amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect RabbitMQ", "Connected to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a Channel", "Channel opened")

	err = ch.ExchangeDeclare(
		commandsExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to open a Channel",  commandsExchange + " exchange declared")
	receiverQueue = QueueDeclare("", true)
	QueueBind(receiverQueue.Name, commandsExchange)

	err = ch.ExchangeDeclare(
		onlineExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to open a Channel",  onlineExchange + " exchange declared")
	onlineQueue = QueueDeclare("", true)
	QueueBind(onlineQueue.Name, onlineExchange)

	return
}

func CloseConnectionAndChannel() (err error){
	err = ch.Close()
	failOnError(err, "Failed to close channel", "Channel closed")
	err = conn.Close()
	failOnError(err, "Failed to close connection", "Connection closed")
	return
}


func QueueDeclare(queueName string, autoDelete bool) (queue amqp.Queue) {
	queue, err := ch.QueueDeclare(
		queueName,
		false,
		autoDelete,
		false,
		false,
		nil,
		)
	failOnError(err, "Failed to declare Queue " + queue.Name, queue.Name + " queue declared")
	return
}

func QueueBind(queueName string, exchangeName string){
	err := ch.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil)
	failOnError(err, "Failed to bind " +queueName+ " queue to " +exchangeName+ " exchange.",
		queueName+ " queue binded to " +exchangeName+ " exchange")
}

func Consumer(queueName string) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,)
	failOnError(err, "Failed to register consumer to " +queueName+ " queue", "")
	return msgs
}

func Publisher(exchangeName string, key string, publishing amqp.Publishing){
	err := ch.Publish(
		exchangeName,
		key,
		false,
		false,
		publishing)
	failOnError(err, "Failed to publish message", "")
}



func CommandsExchange() string{ return commandsExchange}

func ReceiverQueue() string{ return receiverQueue.Name}

func OnlineExchange() string{ return onlineExchange}

func OnlineQueue() string{ return onlineQueue.Name}

