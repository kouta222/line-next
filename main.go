package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("Ingress established at:", listener.URL())

	return http.Serve(listener, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	Secret := os.Getenv("SECRET")
	Token := os.Getenv("TOKEN")

	bot, err := linebot.New(Secret, Token)
	if err != nil {
		log.Printf("Error creating bot client: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Printf("Invalid signature error: %v", err)
			http.Error(w, "Invalid signature", http.StatusBadRequest)
		} else {
			log.Printf("Parse request error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			message := event.Message.(*linebot.TextMessage).Text + "とおしゃっています"
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
				log.Printf("reply message err: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
}
