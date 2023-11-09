package publishers

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/tes-rmq/configs"
	"example.com/tes-rmq/models"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

func Register(ctx *gin.Context) {
	var user models.User

	mqConn, err := configs.ConnectRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	// defer mqConn.Close()

	channel, err := configs.ChannelRabbitMQ(mqConn)
	if err != nil {
		log.Fatal(err)
	}
	// defer channel.Close()

	err = configs.ExchangeDeclareRabbitMQ(channel, "registration_exchange", "direct")
	if err != nil {
		log.Fatal(err)
	}

	err = configs.QueueBindRabbitMQ(channel, "registration_queue", "registration", "registration_exchange")
	if err != nil {
		log.Fatal(err)
	}

	registrationData, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = channel.PublishWithContext(ctx,
		"registration_exchange",
		"registration",
		false,
		false,
		amqp091.Publishing {
			ContentType: "application/json",
			Body: registrationData,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send registration data to RabbitMQ"})
		return
	}
}