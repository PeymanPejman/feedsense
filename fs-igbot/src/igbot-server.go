package main

import (
	"encoding/json"
	pb "feedsense/fs-igbot/src/protos"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port       = ":30200"
	IG_API_URL = "https://api.instagram.com"
)

type Posts struct {
	Pagination struct {
	} `json:"pagination"`
	Data []struct {
		ID   string `json:"id"`
		User struct {
			ID             string `json:"id"`
			FullName       string `json:"full_name"`
			ProfilePicture string `json:"profile_picture"`
			Username       string `json:"username"`
		} `json:"user"`
		Images struct {
			Thumbnail struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"thumbnail"`
			LowResolution struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"low_resolution"`
			StandardResolution struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"standard_resolution"`
		} `json:"images"`
		CreatedTime string `json:"created_time"`
		Caption     struct {
			ID          string `json:"id"`
			Text        string `json:"text"`
			CreatedTime string `json:"created_time"`
			From        struct {
				ID             string `json:"id"`
				FullName       string `json:"full_name"`
				ProfilePicture string `json:"profile_picture"`
				Username       string `json:"username"`
			} `json:"from"`
		} `json:"caption"`
		UserHasLiked bool `json:"user_has_liked"`
		Likes        struct {
			Count int `json:"count"`
		} `json:"likes"`
		Tags     []interface{} `json:"tags"`
		Filter   string        `json:"filter"`
		Comments struct {
			Count int `json:"count"`
		} `json:"comments"`
		Type         string        `json:"type"`
		Link         string        `json:"link"`
		Location     interface{}   `json:"location"`
		Attribution  interface{}   `json:"attribution"`
		UsersInPhoto []interface{} `json:"users_in_photo"`
	} `json:"data"`
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
}

type Comments struct {
	Data []struct {
		CreatedTime string `json:"created_time"`
		Text        string `json:"text"`
		From        struct {
			Username       string `json:"username"`
			ProfilePicture string `json:"profile_picture"`
			ID             string `json:"id"`
			FullName       string `json:"full_name"`
		} `json:"from"`
		ID string `json:"id"`
	} `json:"data"`
}

// server is used to implement AppBot
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) GetClientThreads(ctx context.Context, in *pb.Client) (*pb.GetClientThreadsResponse, error) {
	if in.User.AccessToken == "" {
		return nil, nil
	}
	str := []string{IG_API_URL, "/v1/users/self/media/recent/?&access_token=", in.User.AccessToken}
	url := strings.Join(str, "")

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

	ts := []*pb.Thread{}

	for _, p := range posts.Data {
		ts = append(ts, &pb.Thread{
			Id:    p.ID,
			Title: p.Caption.Text,
			Owner: &pb.Client{
				User: &pb.User{
					UserId:      p.User.Username,
					AccessToken: in.User.AccessToken},
				Agent: pb.Client_INSTAGRAM},
			CommentCount: int32(p.Comments.Count),
			Type:         pb.Thread_IMAGE,
			Url:          p.Link,
		})
	}
	return &pb.GetClientThreadsResponse{Threads: ts}, nil
}

func (s *server) GetThreadComments(ctx context.Context, th *pb.Thread) (*pb.GetThreadCommentsResponse, error) {
	if th.Owner.User.AccessToken == "" {
		log.Fatalf("Not authenticated")
		return nil, nil
	}

	if th.Id == "" {
		log.Fatalf("Unknown thread")
		return nil, nil
	}

	comments := Comments{}

	url := fmt.Sprintf("https://api.instagram.com/v1/media/%v/comments?access_token=%v", th.Id, th.Owner.User.AccessToken)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &comments)
	if err != nil {
		fmt.Println(err)
	}

	result := []*pb.Comment{}
	for _, comment := range comments.Data {
		cm := &pb.Comment{
			CommentId: comment.ID,
			Text:      comment.Text,
			Owner:     comment.From.Username,
			//CreatedTime: comment.CreatedTime,
			Thread: th,
		}
		result = append(result, cm)
	}

	return &pb.GetThreadCommentsResponse{Comments: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAppBotServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
