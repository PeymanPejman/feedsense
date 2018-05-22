package main

import (
	"context"
	pb "feedsense/fs-fe/src/protos"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

const (
	BaseURL = "http://107.178.250.178"
	Address = "fs-igbot:fs-igbot-rpc"
	Name    = "fs-fe"
)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/IGLoginCallback", IGLoginCallback)
	http.HandleFunc("/posts/", ThreadsHandler)
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

// PostsHnadler handles the fetching and showing threads
func ThreadsHandler(w http.ResponseWriter, r *http.Request) {
	cl := &pb.Client{
		User: &pb.User{
			UserId: r.FormValue("userid"), AccessToken: r.FormValue("access_token")},
		Agent: 0}

	conn, err := grpc.Dial(Address, grpc.WithInsecure())
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
	ShowThreads(w, t)
}

// ShowThreads serves info (post id, post caption, time created, comment count) about recenet threads
func ShowThreads(w http.ResponseWriter, t *pb.GetClientThreadsResponse) {
	temp, err := template.ParseFiles("views/posts.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	p := PostsContent{Posts: PrepThreads(t)}
	err = temp.Execute(w, p)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// PrepThreads prepares necessary properties of a thread to show to user
func PrepThreads(t *pb.GetClientThreadsResponse) map[string][]string {
	res := make(map[string][]string)

	for _, thread := range t.Threads {
		attr := []string{thread.Title, thread.CreatedTime.String(), fmt.Sprint(thread.CommentCount)}
		res[thread.Id] = attr
	}

	return res
}

/******** HTTP Authn callbacks *********/

// IGLoginCallback creates a user with credentials from IG
func IGLoginCallback(w http.ResponseWriter, r *http.Request) {
	user := IGAuthCred{}
	IGLogin(w, r, &user)
	url := fmt.Sprintf("%v/posts/?access_token=%v&userid=%v", BaseURL, user.Token, user.User.Id)
	http.Redirect(w, r, url, 302)
}
