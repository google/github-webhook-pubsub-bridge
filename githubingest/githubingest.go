// Package githubingest defines a framework for processing events from GitHub webhooks
package githubingest

import (
	"context"
	"net/http"
	"time"

	"github.com/google/logger"
	"github.com/google/github-webhook-pubsub-bridge/event"
)

type processor interface {
	Process(ctx context.Context, ev *event.Event) error
}

// Server defines a handler for accepting GitHub web hooks and processing them.
type Server struct {
	Processor processor
}

// Handler is the main entrypoint for the Github event webhook.
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("request data: %+v", *r)

	event, err := event.ParseFromRequest(r, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("Error parsing payload: %v\n", err)
		return
	}

	err = s.Processor.Process(r.Context(), event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Errorf("Error processing event: %v", err)
	}
}
