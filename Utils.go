package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getCommentFullName(shareLink string) (string, error) {

	fmt.Println("share link: " + shareLink)
	fullNameIndex := strings.Index(shareLink, "comment/")

	if fullNameIndex == -1 {
		return "", errors.New(" 'comment' not found in share link path")
	}

	fullNameIndex += 8

	if len(shareLink) < fullNameIndex+7 {
		return "", errors.New("share link was cut off too soon")
	}

	return "t1_" + shareLink[fullNameIndex:fullNameIndex+7], nil

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
	req.Header.Add("User-Agent", "Vote Buddy 1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		return &replyCommentData, errors.New("Error sending/receiving api/comment/ request: " + err.Error())
	}

	//fmt.Println(res.Status)

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

	//fmt.Println(replyCommentData)

	return &replyCommentData, nil
}
