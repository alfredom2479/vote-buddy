package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	if err := commentInfo.getCommentInfo(&httpClient, redditToken, "t1_l9nb9df", ApiDomain); err != nil {
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

	generatedReplyComment, err := generateReply(&httpClient, openAIClient, contentForOpenAIMessage)
	if err != nil {
		log.Fatal("Error generating mean reply comment" + err.Error())
	}

	err = sendReply(&httpClient, "t1_l9nb9df", generatedReplyComment, redditToken)
	if err != nil {
		log.Fatal("Error sending generated reply comment to reddit")
	}
	fmt.Println("Comment Succesfully replied to")
}

func sendReply(httpClient *http.Client, parentComment, replyBody, accessToken string) error {

	formData := url.Values{}
	formData.Set("api_type", "json")
	formData.Set("return_rtjson", "true")
	formData.Set("parent", parentComment)
	formData.Set("text", replyBody)

	req, err := http.NewRequest("POST", ApiDomain+"/api/comment", strings.NewReader(formData.Encode()))
	if err != nil {
		return errors.New("Error making new HTTP request to /api/comment: " + err.Error())
	}

	req.Header.Add("AUthorization", "bearer "+accessToken)
	//req.Header.Add("User-Agent") I dont think i need this

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.New("Error sending/receiving api/comment/ request: " + err.Error())
	}

	fmt.Println(res.Status)

	if res.StatusCode != 200 {
		return errors.New("HTTP response not OK: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error reading response body: " + err.Error())
	}

	fmt.Println(string(body))

	return nil
}
