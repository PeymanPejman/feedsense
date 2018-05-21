package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const Base_Url = "http://35.233.157.251:8080"

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/IGLoginCallback", IGLoginCallback)
	http.HandleFunc("/posts", PostsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/******** HTTP Endpoints *********/

// HomePage serves the first page users land on
func HomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/home.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, nil)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

//
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("access_token") == "" || r.FormValue("userid") == "" {
		log.Print("User parameteres not provided")
		return
	}
	user := IGAuthCred{Token: r.FormValue("access_token")}
	user.User.Id = r.FormValue("userid")
	ShowPosts(w, r, &user)
}

//ShowPosts serves info (post id, post caption, time created, comment count) about recenet posts
func ShowPosts(w http.ResponseWriter, r *http.Request, user *IGAuthCred) {
	posts := GetPosts(w, r, user)
	t, err := template.ParseFiles("views/posts.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	p := PostsContent{Posts: PrepPosts(&posts)}
	err = t.Execute(w, p)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// PrepPosts prepares necessary properties of a post to show to user
func PrepPosts(posts *Posts) map[string][]string {
	res := make(map[string][]string)

	for _, post := range posts.Data {
		attr := []string{post.Caption.Text, post.CreatedTime, strconv.Itoa(post.Comments.Count)}
		res[post.ID] = attr
	}

	return res
}

/******** HTTP Authn callbacks *********/

// IGLoginCallback creates a user with credentials from IG
func IGLoginCallback(w http.ResponseWriter, r *http.Request) {
	user := IGAuthCred{}
	IGLogin(w, r, &user)
	url := fmt.Sprintf("%v/posts/?access_token=%v&userid=%v", Base_Url, user.Token, user.User.Id)
	http.Redirect(w, r, url, 302)
}
