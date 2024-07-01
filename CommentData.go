package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (commentData *CommentResponseData) getCommentInfo(httpClient *http.Client, accessToken string, commentFullName string, apiUrl string) error {

	params := url.Values{}

	params.Add("id", commentFullName)

	commentInfoEndpoint := apiUrl + "/api/info?" + params.Encode()

	req, err := http.NewRequest("GET", commentInfoEndpoint, nil)
	if err != nil {
		return errors.New("Error making new HTTP request to comment info endpoint" + err.Error())
	}

	req.Header.Add("Authorization", "bearer "+accessToken)
	req.Header.Add("User-Agent", "Vote Buddy 1.0")

	fmt.Println("sneding req!")
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

func (commentData *CommentResponseData) createContentString(position string) (string, error) {

	contentString := ""
	commentMainData := &commentData.Data.Children[0]

	if commentMainData.Data.Body == "" || commentMainData.Data.Subreddit == "" {
		return contentString, errors.New("comment Body or subreddit not found")
	}

	//working on this
	contentString += "comment thread: '" + commentMainData.Data.Author + "-" + commentMainData.Data.Body + "', subreddit: '" +
		commentMainData.Data.Subreddit + "', position: '" + position + "'"

	return contentString, nil

}
