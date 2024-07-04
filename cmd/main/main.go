package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alfredom2479/vote-buddy/internal/reddit"
	"github.com/alfredom2479/vote-buddy/internal/server"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redditUsername := os.Getenv("REDDIT_USERNAME")
	redditPassword := os.Getenv("REDDIT_PW")
	redditClientId := os.Getenv("REDDIT_CLIENT_ID")
	redditClientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	redditToken := os.Getenv("REDDIT_TOKEN")

	if redditToken == "" {
		httpClient := http.Client{}

		fmt.Println("Getting new auth token...")
		if redditUsername == "" || redditPassword == "" || redditClientId == "" || redditClientSecret == "" {
			log.Fatal("reddit auth data is missing from environment variables")
		}
		redditAuthTokenData := reddit.CreateRedditAuthTokenData()
		if err := redditAuthTokenData.GetAuthToken(&httpClient, redditUsername, redditPassword, redditClientId, redditClientSecret); err != nil {
			log.Fatal("Error getting reddit auth token data: " + err.Error())
		}

		fmt.Println(redditAuthTokenData)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/", server.HandleGetRoot)
	mux.HandleFunc("/submit-link", server.HandleCommentLinkSubmission)

	fmt.Println("sever is listening")
	http.ListenAndServe(":1234", mux)

}

func createRedditAuthTokenData() {
	panic("unimplemented")
}
