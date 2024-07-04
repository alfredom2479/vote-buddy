package reddit

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const API_AUTH_DOMAIN = "https://www.reddit.com/api/v1/access_token?scope=*"

type RedditAuthTokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func CreateRedditAuthTokenData() RedditAuthTokenData {
	return RedditAuthTokenData{
		AccessToken: "",
		TokenType:   "",
		ExpiresIn:   0,
		Scope:       "",
	}
}

func (tokenData *RedditAuthTokenData) GetAuthToken(httpClient *http.Client, username string, password string, clientID string, clientSecret string) error {

	data := url.Values{}
	data.Add("grant_type", "password")
	data.Add("username", username)
	data.Add("password", password)

	req, err := http.NewRequest("POST", API_AUTH_DOMAIN, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.New("Error making new reddit auth reqeust: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Vote Buddy 1.0")
	req.SetBasicAuth(clientID, clientSecret)

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.New("Error sending/receiving http response: " + err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error reading response body: " + err.Error())
	}

	if err := json.Unmarshal(body, &tokenData); err != nil {
		return errors.New("Error unmarshaling response body into RedditAuthToken object: " + err.Error())
	}

	return nil
}
