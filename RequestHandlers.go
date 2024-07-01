package main

import (
	"fmt"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func handleGetRoot(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func handleCommentLinkSubmission(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		formVals := r.Form
		//go func() {

		commentFullName, err := getCommentFullName(formVals["share-link"][0])
		if err != nil {
			fmt.Println("Error getting comment full name: " + err.Error())
			fmt.Fprintf(w, "<p>Comment Share Link is not valid</p>")
			return
		}
		fmt.Println("comment full name:" + commentFullName)

		//Start the comment reply process

		fmt.Println(isGeneratingComment.Load())
		if isGeneratingComment.Load() {
			fmt.Println("Already generating comment")
			fmt.Fprintf(w, "<p>Comment already being generated</p>")
			return
		}

		isGeneratingComment.Store(true)

		httpClient := http.Client{}
		redditToken := os.Getenv("REDDIT_TOKEN")

		var commentInfo CommentResponseData

		if err := commentInfo.getCommentInfo(&httpClient, redditToken, commentFullName, ApiDomain); err != nil {
			fmt.Println("Error getting comment info: " + err.Error())
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			isGeneratingComment.Store(false)
			return
		}

		fmt.Println("parent id: " + commentInfo.Data.Children[0].Data.ParentID)
		fmt.Println("author name" + commentInfo.Data.Children[0].Data.Author)

		contentForOpenAIMessage, err := commentInfo.createContentString("agree")
		if err != nil {
			fmt.Println("Error creating content string " + err.Error())
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			isGeneratingComment.Store(false)
			return
		}
		fmt.Println("content for AI: " + contentForOpenAIMessage)

		openAIToken := os.Getenv("OPENAI_TOKEN")
		if openAIToken == "" {
			fmt.Println("Error finding openai api token")
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			isGeneratingComment.Store(false)
			return
		}

		openAIClient := openai.NewClient(openAIToken)

		generatedReplyComment, err := generateReply(&httpClient, openAIClient, contentForOpenAIMessage)
		if err != nil {
			fmt.Println("Error generating mean reply comment" + err.Error())
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			isGeneratingComment.Store(false)
			return
		}

		replyData, err := sendReply(&httpClient, commentFullName, generatedReplyComment, redditToken)
		if err != nil {
			fmt.Println("Error sending generated reply comment to reddit")
			fmt.Fprintf(w, "<p>Could not generate comment</p>")
			isGeneratingComment.Store(false)
			return
		}

		isGeneratingComment.Store(false)
		//fmt.Println(replyData)
		fmt.Println("Comment Succesfully replied to")
		fmt.Fprintf(w, "<p>Vote Buddy Has Replied</p> <a href='https://www.reddit.com"+replyData.Permalink+"' class='underline text-blue-800'>See Generated Reply</a>")
		//}()
	}
}
