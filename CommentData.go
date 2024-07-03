package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type CommentResponseData struct {
	Kind string `json:"kind"`
	Data struct {
		After    string `json:"after"`
		Dist     int    `json:"dist"`
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				SubredditID    string  `json:"subreddit_id"`
				LinkTitle      string  `json:"link_title"`
				Subreddit      string  `json:"subreddit"`
				Title          string  `json:"title"`
				Selftext       string  `json:"selftext"`
				LinkAuthor     string  `json:"link_author"`
				Author         string  `json:"author"`
				ParentID       string  `json:"parent_id"`
				AuthorFullname string  `json:"parent_fullname"`
				Body           string  `json:"body"`
				BodyHTML       string  `json:"body_html"`
				LinkID         string  `json:"link_id"`
				Permalink      string  `json:"permalink"`
				LinkPermalink  string  `json:"link_permalink"`
				Name           string  `json:"name"`
				CreatedUTC     float64 `json:"created_utc"`
				LinkURL        string  `json:"link_url"`
				Replies        string  `json:"replies"`
			} `json:"data"`
		} `json:"children"`
		Before string `json:"before"`
	} `json:"data"`
}

type ReplyCommentData struct {
	SubredditID string `json:"subreddit_id"`
	Subreddit   string `json:"subreddit"`
	Name        string `json:"name"`
	LinkID      string `json:"link_id"`
	ID          string `json:"id"`
	Author      string `json:"author"`
	Body        string `json:"body"`
	Permalink   string `json:"permalink"`
}

const myUsername = "No-Atmosphere9068"

func (commentData *CommentResponseData) getCommentInfo(httpClient *http.Client, accessToken, commentFullName, apiUrl string) error {

	params := url.Values{}

	params.Add("id", commentFullName)

	commentInfoEndpoint := apiUrl + "/api/info?" + params.Encode()

	req, err := http.NewRequest("GET", commentInfoEndpoint, nil)
	if err != nil {
		return errors.New("Error making new HTTP request to comment info endpoint" + err.Error())
	}

	req.Header.Add("Authorization", "bearer "+accessToken)
	req.Header.Add("User-Agent", "Vote Buddy 1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.New("Error sending/receiving http comment req/res: " + err.Error())
	}

	if res.StatusCode != 200 {
		return errors.New("HTTP code not 200: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error reading response body: " + err.Error())
	}

	if err := json.Unmarshal([]byte(string(body)), &commentData); err != nil {
		return errors.New("Erro unmarshaling into commentData struct: " + err.Error())
	}

	return nil
}

func createContentString(commentDataSlice []CommentResponseData, postData *CommentResponseData, position string) (string, error) {

	mainPostData := postData.Data.Children[0].Data

	if len(commentDataSlice) < 1 {
		return "", errors.New("commentDataSlice is empty")
	}

	contentString := ""
	commentAuthor := ""

	for _, commentData := range commentDataSlice {

		commentMainData := commentData.Data.Children[0]

		if commentMainData.Data.Body == "" || commentMainData.Data.Subreddit == "" {
			return contentString, errors.New("comment Body or subreddit not found")
		}

		commentAuthor = commentMainData.Data.Author

		if commentAuthor == myUsername {
			commentAuthor = "(ME)"
		}

		contentString = commentAuthor + "-\"" +
			commentMainData.Data.Body + "\"\n" + contentString

	}

	contentString = "Comment thread: \n" + contentString

	contentString += ",subreddit: '" + commentDataSlice[0].Data.Children[0].Data.Subreddit +
		"',\nposition: '" + position +
		"',\npost title: '" + mainPostData.Title + "'"

	if mainPostData.Selftext != "" {
		contentString += ",\npost body text: '" + mainPostData.Selftext + "'"
	}
	return contentString, nil

}
