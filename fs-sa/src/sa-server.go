package main

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	// Imports the Google Cloud Natural Language API client package.
	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "feedsense/fs-sa/src/protos"
)

const (
	BaseURL       = "http://107.178.250.178"
	Address_igbot = "107.178.250.178:30200"
	Name          = "fs-fe"
	port          = ":30100"
)

// SentimentResponse to hold the result coming from NLP calls
type SentimentResponse struct {
	Comment   string
	Sentiment string
	Score     float32
}

// server is used to implement SentimentAnalysis Service
type server struct{}

// GetSentiment gets sentiment analysis response from Google Cloud NLP
func GetSentiment(wg *sync.WaitGroup, c chan SentimentResponse, text string) {
	defer wg.Done()

	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the text to analyze.

	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	//return sentiment
	if sentiment.DocumentSentiment.Score >= 0 {
		response := SentimentResponse{Comment: text, Sentiment: "Positive : " + strconv.FormatFloat(float64(sentiment.DocumentSentiment.Score), 'f', 3, 32),
			Score: sentiment.DocumentSentiment.Score}
		c <- response
	} else {
		response := SentimentResponse{Comment: text, Sentiment: "Negative : " + strconv.FormatFloat(float64(sentiment.DocumentSentiment.Score), 'f', 3, 32),
			Score: sentiment.DocumentSentiment.Score}
		c <- response
	}
}

// GetThreadComments makes an RPC call to *Bot to get all comments in a thread
func GetThreadComments(th *pb.Thread) (*pb.GetThreadCommentsResponse, error) {
	conn, err := grpc.Dial(Address_igbot, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAppBotClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	t, err := c.GetThreadComments(ctx, th)
	if err != nil {
		log.Fatalf("could not get threads: %v", err)
	}

	return t, err
}

func (s *server) GetCurrentThreadSentiment(tx context.Context, in *pb.Thread) (*pb.ThreadSentiment, error) {
	var wg sync.WaitGroup

	comments, err := GetThreadComments(in)

	if err != nil {
		log.Fatalf("Could not get comments: %v", err)
		return nil, nil
	}

	sentiments := make(map[string]string)

	queue := make(chan SentimentResponse, 30)
	for _, cm := range comments.Comments {
		wg.Add(1)
		go GetSentiment(&wg, queue, cm.Text)
	}

	wg.Wait()
	close(queue)
	sum := float32(0.0)
	for q := range queue {
		sentiments[q.Comment] = q.Sentiment
		sum += q.Score
	}
	sentiment := &pb.ThreadSentiment{
		Score:  float64(sum) / float64(len(comments.Comments)),
		Thread: in}
	return sentiment, nil
}

func (s *server) GetAllCurrentClientThreadsSentiment(tx context.Context, cl *pb.Client) (*pb.GetAllCurrentClientThreadsSentimentResponse, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSentimentAnalysisServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
