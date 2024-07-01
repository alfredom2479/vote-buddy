package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
)

//for future use,

const ApiAuthDomain = "https://www.reddit.com/api/v1/access_token?scope=*"
const ApiDomain = "https://oauth.reddit.com"

var isGeneratingComment atomic.Bool

func main() {

	isGeneratingComment.Store(false)
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
		redditAuthTokenData := createRedditAuthTokenData()
		if err := redditAuthTokenData.getAuthToken(&httpClient, redditUsername, redditPassword, redditClientId, redditClientSecret); err != nil {
			log.Fatal("Error getting reddit auth token data: " + err.Error())
		}

		fmt.Println(redditAuthTokenData)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/", handleGetRoot)
	mux.HandleFunc("/submit-link", handleCommentLinkSubmission)

	fmt.Println("sever is listening")
	http.ListenAndServe(":1234", mux)

}
