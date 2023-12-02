package githubingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/pubsub"
	        "github.com/google/github-webhook-pubsub-bridge/event"
)

var errPublish = errors.New("error publishing message")

// Topic is just a small wrapper around pubsub.Topic that collapses the
// Publish...Get methods into a single call (makes testing easier, too).
type Topic struct {
	*pubsub.Topic
}

// PublishMessage is a simplified wrapper for pubsub topics - it publishes and
// just returns the error.
func (t *Topic) PublishMessage(ctx context.Context, msg *pubsub.Message) error {
	_, err := t.Publish(ctx, msg).Get(ctx)
	return err
}

type messagePublisher interface {
	PublishMessage(context.Context, *pubsub.Message) error
}

// PubSubProcessor publishes pubsub events to the given topic for each processed event.
type PubSubProcessor struct {
	// Topic is The full private feed - events from all repositories (including private ones).
	Topic messagePublisher
	// PublicOnlyTopic is a feed of only events from public repositories.
	PublicOnlyTopic messagePublisher
}

// Process runs the pubsub processor for a single event.
func (psp *PubSubProcessor) Process(ctx context.Context, ev *event.Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed marshalling event: %v", err)
	}
	msg := &pubsub.Message{Data: b, Attributes: extractAttributes(ev)}
	if err = psp.Topic.PublishMessage(ctx, msg); err != nil {
		return fmt.Errorf("%w: %v", errPublish, err)
	}

	// Publish to the public-only feed if it's not private
	if ev.RepoInfo != nil && !ev.RepoInfo.IsPrivate {
		if err = psp.PublicOnlyTopic.PublishMessage(ctx, msg); err != nil {
			return fmt.Errorf("%w: %v", errPublish, err)
		}
	}

	return nil
}

func extractAttributes(ev *event.Event) map[string]string {
	a := make(map[string]string)
	if ev == nil {
		return a
	}
	a["type"] = ev.Type

	r := ev.RepoInfo
	if r == nil {
		return a
	}

	if r.OwnerName != "" {
		a["owner"] = r.OwnerName
	}
	if r.Name != "" {
		a["repository"] = r.Name
	}
	if r.OrganizationName != "" {
		a["organization"] = r.OrganizationName
	}

	return a
}
