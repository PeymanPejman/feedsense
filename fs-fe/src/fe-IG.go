package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const IG_API_URL = "https://api.instagram.com"
const REDIRECT_URL = "http://127.0.0.1:8080/IGLoginCallback"
const CLIENT_SECRET = "a8828a03a5144f43a863b2d49d69f60f"
const CLIENT_ID = "2b651acd74af417686fba2086af5b962"

//IGLogin handles the Oauth2 flow
func IGLogin(w http.ResponseWriter, r *http.Request, user *IGAuthCred) {
	if r.FormValue("code") == "" {
		log.Print("Code was not recieved")
		return
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
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Print("Error reading the body", err)
	}

	err = json.Unmarshal(b, &user)

	if err != nil {
		log.Print("Error unmarshalling the reponse", err)
	}

	err = resp.Body.Close()

	if err != nil {
		log.Print("Cannot close the body", err)
	}
}

// GetPosts retrive metadata of all the recent posts in IG and returns Posts
func GetPosts(w http.ResponseWriter, r *http.Request, user *IGAuthCred) Posts {
	//Replace with call to IGBot to get posts
	if user.Token == "" {
		log.Print("User not authenticated")
		http.Redirect(w, r, "/", 301)
	}

	s := []string{IG_API_URL, "/v1/users/self/media/recent/?&access_token=", user.Token}
	url := strings.Join(s, "")

	resp, err := http.Get(url)
	if err != nil {
		log.Print("Could not get posts from IG : ", err)
	}

	posts := Posts{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Cannot close the body", err)
	}

	err = json.Unmarshal(b, &posts)
	if err != nil {
		fmt.Println(err)
	}

	resp.Body.Close()

	return posts

}
