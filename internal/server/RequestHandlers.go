package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alfredom2479/vote-buddy/internal/openai"
	"github.com/alfredom2479/vote-buddy/internal/reddit"
)

const API_DOMAIN = "https://oauth.reddit.com"

func HandleGetRoot(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func HandleCommentLinkSubmission(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		formVals := r.Form

		fmt.Println(formVals)

		voteBuddyPosition := "agree"

		commentFullName, err := reddit.GetCommentFullName(formVals["share-link"][0])
		if err != nil {
			fmt.Println("Error getting comment full name: " + err.Error())
			fmt.Fprintf(w, "<p>Comment Share Link is not valid</p>")
			return
		}

		voteOption := formVals["voteOption"][0]

		switch voteOption {
		case "downvote":
			voteBuddyPosition = "disagree"
		case "upvote":
			voteBuddyPosition = "agree"
		default:
			voteBuddyPosition = "agree"

		}

		httpClient := http.Client{}
		redditToken := os.Getenv("REDDIT_TOKEN")

		var commentInfoSlice []reddit.CommentResponseData
		var postInfo reddit.CommentResponseData

		currContentName := commentFullName
		numOfReqs := 0

		for {

			commentInfo := reddit.CommentResponseData{}

			numOfReqs += 1

			if err := commentInfo.GetCommentInfo(&httpClient, redditToken, currContentName, API_DOMAIN); err != nil {
				fmt.Println("Error getting comment info: " + err.Error())
				fmt.Fprintf(w, "<p>Could not generate comment</p>")
				return
			}

			commentInfoSlice = append(commentInfoSlice, commentInfo)

			currContentName = commentInfo.Data.Children[0].Data.ParentID

			if !strings.HasPrefix(currContentName, "t1") || numOfReqs > 15 {

				if err = postInfo.GetCommentInfo(&httpClient, redditToken, currContentName, API_DOMAIN); err != nil {
					fmt.Println("Error getting comment info: " + err.Error())
					fmt.Fprintf(w, "<p>Could not generate comment</p>")
					return
				}

				break
			}
		}

		contentForOpenAIMessage, err := reddit.CreateContentString(commentInfoSlice, &postInfo, voteBuddyPosition)
		if err != nil {
			fmt.Println("Error creating content string " + err.Error())
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			return
		}
		fmt.Println("content for AI: " + contentForOpenAIMessage)

		/*
			openAIToken := os.Getenv("OPENAI_TOKEN")
			if openAIToken == "" {
				fmt.Println("Error finding openai api token")
				fmt.Fprintf(w, "<p>Could not generate comment</p>")
				return
			}

			openAIClient := openai.NewClient(openAIToken)
		*/

		generatedReplyComment, err := openai.GenerateReply(&httpClient, contentForOpenAIMessage)
		if err != nil {
			fmt.Println("Error generating mean reply comment" + err.Error())
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			return
		}

		replyData, err := reddit.SendReply(&httpClient, commentFullName, generatedReplyComment, redditToken)
		if err != nil {
			fmt.Println("Error sending generated reply comment to reddit")
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			return
		}
		if replyData.ID == "" {
			fmt.Fprintf(w, "<p>Could not post reply (original post may be deleted)</p>")
			return
		}

		fmt.Println("Comment Succesfully replied to")
		fmt.Fprintf(w, "<p>Vote Buddy Has Replied</p> <a target='_blank' href='https://www.reddit.com"+replyData.Permalink+"' class='underline text-blue-800'>See Generated Reply</a>")
	}
}
