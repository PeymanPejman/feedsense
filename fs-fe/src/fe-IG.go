package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const IG_API_URL = "https://api.instagram.com"
const REDIRECT_URL = BaseURL + "/IGLoginCallback"
const CLIENT_SECRET = "a8828a03a5144f43a863b2d49d69f60f"
const CLIENT_ID = "2b651acd74af417686fba2086af5b962"

//IGLogin handles the Oauth2 flow
func IGLogin(w http.ResponseWriter, r *http.Request, user *IGAuthCred) bool {
	if r.FormValue("code") == "" {
		log.Print("Code was not recieved")
		return false
	}
	apiURL := IG_API_URL
	resource := "/oauth/access_token"
	data := url.Values{}
	data.Add("code", r.FormValue("code"))
	data.Add("redirect_uri", REDIRECT_URL)
	data.Add("grant_type", "authorization_code")
	data.Add("client_secret", CLIENT_SECRET)
	data.Add("client_id", CLIENT_ID)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))

	if err != nil {
		log.Print("Eorror creating the POST request : ", err)
		return false
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Print("Error reading the body", err)
		return false
	}

	err = json.Unmarshal(b, &user)

	if err != nil {
		log.Print("Error unmarshalling the reponse", err)
		return false
	}

	err = resp.Body.Close()

	if err != nil {
		log.Print("Cannot close the body", err)
		return false
	}

	return true
}

//
func IsAccessTokenAlive(token string) bool {
	str := []string{IG_API_URL, "/v1/users/self/?&access_token=", token}
	url := strings.Join(str, "")

	resp, _ := http.Get(url)
	return resp.StatusCode == http.StatusOK
}
