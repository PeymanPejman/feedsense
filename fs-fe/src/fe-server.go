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
	BaseURL       = "http://127.0.0.1"
	Address_igbot = "http://127.0.0.1:30200"
	Address_sa    = "http://127.0.0.1:30100"
	Name          = "fs-fe"
)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/IGLoginCallback", IGLoginCallback)
	http.HandleFunc("/threads/", ThreadsHandler)
	http.HandleFunc("/sentiment/", SentimentHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
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

// SentimentHandler handles sentiment analysis aquisition
func SentimentHandler(w http.ResponseWriter, r *http.Request) {
	th := &pb.Thread{
		Id: r.FormValue("threadid"),
		Owner: &pb.Client{
			User: &pb.User{UserId: r.FormValue("userid"),
				AccessToken: r.FormValue("access_token")},
			Agent: pb.Client_INSTAGRAM},
	}

	//TODO: Factor out grpc dial and make a helper

	conn, err := grpc.Dial(Address_sa, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSentimentAnalysisClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cms, err := c.GetCurrentThreadSentiment(ctx, th)
	if err != nil {
		log.Fatalf("could not get threads: %v", err)
	}
	ShowSentiment(w, cms)
}

// ThreadsHnadler handles the fetching and showing threads
func ThreadsHandler(w http.ResponseWriter, r *http.Request) {
	cl := &pb.Client{
		User: &pb.User{
			UserId: r.FormValue("userid"), AccessToken: r.FormValue("access_token")},
		Agent: pb.Client_INSTAGRAM}

	conn, err := grpc.Dial(Address_igbot, grpc.WithInsecure())
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
	ShowThreads(w, t)
}

/******** Helpers *********/

// ShowThreads serves info about recenet threads
func ShowThreads(w http.ResponseWriter, t *pb.GetClientThreadsResponse) {
	temp, err := template.ParseFiles("views/threads.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	p := ThreadsContent{Threads: PrepThreads(t)}
	err = temp.Execute(w, p)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// ShowSentiment serves sentiment info about a given Thread
func ShowSentiment(w http.ResponseWriter, t *pb.ThreadSentiment) {
	temp, err := template.ParseFiles("views/sentiment.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	sent := []string{t.Thread.Id, strconv.FormatFloat(t.Score, 'f', -1, 64)}
	p := SentimentContent{Sentiment: sent}
	err = temp.Execute(w, p)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// PrepThreads prepares necessary properties of a thread to show to user
func PrepThreads(t *pb.GetClientThreadsResponse) map[string][]string {
	res := make(map[string][]string)

	for _, thread := range t.Threads {
		attr := []string{thread.Id, thread.Title, fmt.Sprint(thread.CommentCount)}
		res[thread.Id] = attr
	}

	return res
}

/******** HTTP Authn callbacks *********/

// IGLoginCallback creates a user with credentials from IG
func IGLoginCallback(w http.ResponseWriter, r *http.Request) {
	user := IGAuthCred{}
	IGLogin(w, r, &user)
	url := fmt.Sprintf("%v/threads/?access_token=%v&userid=%v", BaseURL, user.Token, user.User.Id)
	http.Redirect(w, r, url, 302)
}
