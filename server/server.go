package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	red "github.com/ahmetsabri/localhook/redis"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var reponseFromWebhook chan string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	reponseFromWebhook = make(chan string)
	redisClient = red.CreateRedisClient()
	r := mux.NewRouter()
	staticDir := "./static/"

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"index.html")
	}).Methods("GET")

	r.HandleFunc("/connect/{apikey}", connect)
	r.HandleFunc("/webhook/{apikey}", webhook).Methods(http.MethodPost)
	serverPort := os.Getenv("SERVER_PORT")
	color.Cyan("server listenin on 127.0.0.1%s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, r))
}

func connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apikey := vars["apikey"]

	red.SetClientConnected(redisClient, apikey, 1)
	startSse(w, r, apikey)
}

func startSse(w http.ResponseWriter, r *http.Request, apikey string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second * 3)
	defer t.Stop()
	for {
		select {
		case <-clientGone:
			fmt.Printf("Client %s disconnected\n", apikey)
			closeClientConnection(redisClient, apikey)
			return
		case data := <-reponseFromWebhook:
			color.Cyan("Response : %s", data)
			_, err := fmt.Fprintf(w, "%s\n\n", data)
			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
		}
	}
}

func closeClientConnection(client *redis.Client, key string) error {
	if red.CheckClientConnection(client, key) {
		ctx := context.Background()
		err := client.Set(ctx, key, "0", 0).Err()
		if err != nil {
			return fmt.Errorf("could not close client in Redis: %v", err)
		}
		color.Magenta("Client %s connection closed successfully !\n", key)
	}
	return nil

}

func webhook(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	apikey := vars["apikey"]
	ctx := context.Background()

	if redisClient.Get(ctx, apikey).Val() == "0" {
		color.Red("User in not connected")
		fmt.Fprintf(w, "Client is not connected")

		return
	}
	reponseFromWebhook <- string(body)
}
