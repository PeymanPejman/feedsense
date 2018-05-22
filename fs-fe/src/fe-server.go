package main

import (
	"context"
	pb "feedsense/fs-fe/src/protos"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	BaseUrl = "http://35.230.43.0:8080"
	address = "localhost:34000"
	name    = "fs-fe"
)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/IGLoginCallback", IGLoginCallback)
	http.HandleFunc("/posts/", PostsHandler)
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

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	cl := &pb.Client{
		User: &pb.User{
			UserId: r.FormValue("userid"), AccessToken: r.FormValue("access_token")},
		Agent: 0}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAppBotClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	t, err := c.GetClientThreads(ctx, cl)
	if err != nil {
		log.Fatalf("could not get threads: %v", err)
	}
	log.Printf(t.Threads[0].Id)
}

//ShowPosts serves info (post id, post caption, time created, comment count) about recenet posts
/*func ShowPosts(w http.ResponseWriter, r *http.Request, user *IGAuthCred) {
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
}*/

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
	url := fmt.Sprintf("%v/posts/?access_token=%v&userid=%v", BaseUrl, user.Token, user.User.Id)
	http.Redirect(w, r, url, 302)
}
