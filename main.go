// Package main runs the github event ingester that stores events in Spanner and rebroadcasts them over pubsub.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/logger"
	"cloud.google.com/go/pubsub"
	"github.com/google/github-webhook-pubsub-bridge/githubingest"
)

func main() {
	logger.Init("github-webhook-pubsub-bridge", false, false, os.Stdout)
	topic := os.Getenv("TOPIC")
	publicOnlyTopic := os.Getenv("PUBLIC_TOPIC")
	project := os.Getenv("GCP_PROJECT")

	if topic == "" {
		topic = "githubevent"
	}
	if publicOnlyTopic == "" {
		publicOnlyTopic = "githubevent-public"
	}

	cli, err := pubsub.NewClient(context.Background(), project)
	if err != nil {
		logger.Fatalf("failed to make pubsub client: %v", err)
		os.Exit(1)
	}

	ps := &githubingest.PubSubProcessor{
		Topic:           &githubingest.Topic{cli.Topic(topic)},
		PublicOnlyTopic: &githubingest.Topic{cli.Topic(publicOnlyTopic)},
	}
	s := githubingest.Server{
		Processor: ps,
	}
	http.HandleFunc("/", s.Handler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "OK") })
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
