package consumers

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/tes-rmq/configs"
	"example.com/tes-rmq/models"
	"github.com/gin-gonic/gin"
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

	queue, err := configs.QueueDeclareRabbitMQ(channel, "registration_queue")
	if err != nil {
		log.Fatal(err)
	}
	// defer cancel()

	err = configs.QueueBindRabbitMQ(channel, "registration_queue", "registration", "registration_exchange")
	if err != nil {
		log.Fatal(err)
	}

	msg, ok ,err := channel.Get(queue.Name, true)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get a message from the queue"})
			return
		}

		if !ok {
			ctx.JSON(http.StatusNoContent, gin.H{"message": "No message available in the queue"})
			return
		}

		messageBody := string(msg.Body)
		err = json.Unmarshal([]byte(messageBody), &user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal JSON"})
			return
		}
}