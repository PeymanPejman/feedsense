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
	BaseURL       = "http://107.178.250.178/:8080"
	Address_igbot = "107.178.250.178:30200"
	Address_sa    = "107.178.250.178:30100"
	Name          = "fs-fe"
)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/IGLoginCallback", IGLoginCallback)
	http.HandleFunc("/threads/", ThreadsHandler)
	http.HandleFunc("/sentiment/", SentimentHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/******** HTTP Endpoints *********/

// HomePage serves the first page users land on
func HomePage(w http.ResponseWriter, r *http.Request) {
	if IsCookieValid(r) {
		url := fmt.Sprintf("%v/threads/", BaseURL)
		http.Redirect(w, r, url, 302)
	}

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
	if !IsCookieValid(r) {
		url := fmt.Sprintf("%v/", BaseURL)
		http.Redirect(w, r, url, 302)
	}

	th := &pb.Thread{
		Id: r.FormValue("ThreadId"),
		Owner: &pb.Client{
			User:  &pb.User{AccessToken: MustGetValue("access_token", r)},
			Agent: pb.Client_INSTAGRAM},
	}

	//TODO: Factor out grpc dial and make a helper

	conn, err := grpc.Dial(Address_sa, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewSentimentAnalysisClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cms, err := c.GetCurrentThreadSentiment(ctx, th)
	if err != nil {
		log.Fatalf("could not get threads: %v", err)
		return
	}
	ShowSentiment(w, cms)
}

// ThreadsHandler handles the fetching and showing threads
func ThreadsHandler(w http.ResponseWriter, r *http.Request) {
	if !IsCookieValid(r) {
		url := fmt.Sprintf("%v/", BaseURL)
		http.Redirect(w, r, url, 302)
	}

	cl := &pb.Client{
		User: &pb.User{
			AccessToken: MustGetValue("access_token", r)},
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

/******** HTTP Authn callbacks *********/

// IGLoginCallback creates a user with credentials from IG
func IGLoginCallback(w http.ResponseWriter, r *http.Request) {
	user := IGAuthCred{}
	if !IGLogin(w, r, &user) {
		url := fmt.Sprintf("%v/", BaseURL)
		http.Redirect(w, r, url, 302)
	}

	cookie := http.Cookie{Name: "access_token", Value: user.Token}
	http.SetCookie(w, &cookie)
	url := fmt.Sprintf("%v/threads/", BaseURL)
	http.Redirect(w, r, url, 302)
}

/******** Helpers *********/

// IsCookieValid checks to see if the cookie contains access token
func IsCookieValid(r *http.Request) bool {
	c, _ := r.Cookie("access_token")
	if c != nil {
		return IsAccessTokenAlive(MustGetValue("access_token", r))
	}

	return false
}

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

	feel := "Negative"
	if t.Score > 0.0 {
		feel = "Positive"
	}

	sent := []string{t.Thread.Id, strconv.FormatFloat(t.Score, 'f', 2, 64), feel}
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

// MustGetValue returns the value of a Key in a cookie
func MustGetValue(key string, r *http.Request) string {
	c, _ := r.Cookie(key)
	return c.Value
}
