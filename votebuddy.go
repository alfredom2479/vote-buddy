package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	mux.HandleFunc("/", handleGetRoot)
	mux.HandleFunc("/submit-link", handleCommentLinkSubmission)

	fmt.Println("sever is listening")
	http.ListenAndServe(":1234", mux)

}

func sendReply(httpClient *http.Client, parentComment, replyBody, accessToken string) (*ReplyCommentData, error) {

	var replyCommentData ReplyCommentData

	formData := url.Values{}
	formData.Set("api_type", "json")
	formData.Set("return_rtjson", "true")
	formData.Set("parent", parentComment)
	formData.Set("text", replyBody)

	req, err := http.NewRequest("POST", ApiDomain+"/api/comment", strings.NewReader(formData.Encode()))
	if err != nil {
		return &replyCommentData, errors.New("Error making new HTTP request to /api/comment: " + err.Error())
	}

	req.Header.Add("AUthorization", "bearer "+accessToken)
	//req.Header.Add("User-Agent") I dont think i need this

	res, err := httpClient.Do(req)
	if err != nil {
		return &replyCommentData, errors.New("Error sending/receiving api/comment/ request: " + err.Error())
	}

	fmt.Println(res.Status)

	if res.StatusCode != 200 {
		return &replyCommentData, errors.New("HTTP response not OK: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &replyCommentData, errors.New("Error reading response body: " + err.Error())
	}

	if err := json.Unmarshal([]byte(string(body)), &replyCommentData); err != nil {
		return &replyCommentData, errors.New("Error unmarshalling into replyCommentDataStruct" + err.Error())
	}

	fmt.Println(replyCommentData)

	return &replyCommentData, nil
}
