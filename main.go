package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example.com/tes-rmq/configs"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Username string `json:"username"`
}

func main() {
	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := configs.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect(ctx)

	collection := configs.GetCollection(db, "userss")

	mqConn, err := configs.ConnectRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer mqConn.Close()

	channel, err := configs.ChannelRabbitMQ(mqConn)
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	err = configs.ExchangeDeclareRabbitMQ(channel, "registration_exchange", "direct")
	if err != nil {
		log.Fatal(err)
	}

	queue, err := configs.QueueDeclareRabbitMQ(channel, "registration_queue")
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	err = configs.QueueBindRabbitMQ(channel, "registration_queue", "registration", "registration_exchange")
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return	
		}

		err := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&user)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}

		registrationData, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send registration data to RabbitMQ"})
			return
		}

		msg, ok ,err := channel.Get(queue.Name, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get a message from the queue"})
			return
		}

		if !ok {
			c.JSON(http.StatusNoContent, gin.H{"message": "No message available in the queue"})
			return
		}

		messageBody := string(msg.Body)
		err = json.Unmarshal([]byte(messageBody), &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal JSON"})
			return
		}

		_, err = collection.InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": "201",
			"message": "Register berhasil",
			"data": user,
		})

	})

	r.Run(":8000")

}