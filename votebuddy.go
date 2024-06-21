package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

//for future use,

const ApiAuthDomain = "https://www.reddit.com/api/v1/access_token?scope=*"
const ApiDomain = "https://oauth.reddit.com"

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

	httpClient := http.Client{}

	if redditToken == "" {
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

	var commentInfo CommentResponseData

	if err := commentInfo.getCommentInfo(&httpClient, redditToken, "t1_l3hp8d2", ApiDomain); err != nil {
		log.Fatal("Error getting comment info: " + err.Error())
	}

	contentForOpenAIMessage, err := commentInfo.createContentString("agree")
	if err != nil {
		log.Fatal("Error creating content string " + err.Error())
	}
	fmt.Println(contentForOpenAIMessage)

	openAIToken := os.Getenv("OPENAI_TOKEN")
	if openAIToken == "" {
		log.Fatal("Error finding openai api token")
	}

	openAIClient := openai.NewClient(openAIToken)

	if err := generateMeanReply(&httpClient, openAIClient, contentForOpenAIMessage); err != nil {
		log.Fatal("Error generating mean reply comment" + err.Error())
	}

}
