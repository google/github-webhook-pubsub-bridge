// Package event parses github webhooks and converts their data to a manageable Event type.
package event

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
)

const (
	// secretToken is the secret we currently use from GitHub that lets us verify
	// that they actually sent us the request. It's a shared secret between us &
	// GitHub.
	secretToken     = "wkp5H9YShwQkWmV"
	signatureHeader = "X-Hub-Signature"
)

// RepoInfo Contains information about the repository the Event is from.
type RepoInfo struct {
	IsPrivate        bool
	Name             string
	OwnerName        string
	OrganizationName string
}

// Event is a parsed representation of an Event from GitHub.
type Event struct {
	Signature string
	Type      string
	Received  time.Time
	Payload   []byte
	RepoInfo  *RepoInfo `json:"-"` // Do not want to encode this in PubSub
}

// ParseFromRequest returns a parsed Event from the incoming HTTP webhook payload.
func ParseFromRequest(r *http.Request, received time.Time) (*Event, error) {
	payload, err := github.ValidatePayload(r, []byte(secretToken))
	if err != nil {
		return nil, fmt.Errorf("error validating payload: %v", err)
	}
	// NOTE: ValidatePayload ensures the signature header is in the correct form.
	sigParts := strings.Split(r.Header.Get(signatureHeader), "=")

	wht := github.WebHookType(r)
	rep, err := parseRepository(wht, payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %v", err)
	}

	return &Event{
		Signature: sigParts[1],
		Type:      wht,
		Received:  received,
		Payload:   payload,
		RepoInfo:  rep,
	}, nil
}

// String formats the event as a pretty printed string (multi-line).
func (e Event) String() string {
	payload := string(e.Payload)
	if len(payload) > 1000 {
		payload = payload[:1000] + "…"
	}
	return fmt.Sprintf(`{
  Signature: %v,
  Type: %v,
  Received: %v,
  Payload: %v,
	Repository: %v,
}`, e.Signature, e.Type, e.Received, string(payload), e.RepoInfo)
}

func parseRepository(wht string, payload []byte) (*RepoInfo, error) {
	ev, err := github.ParseWebHook(wht, payload)
	if err != nil {
		return nil, err
	}

	// The event field naming is here inconsistent at best.
	switch ev := ev.(type) {
	case *github.OrganizationEvent:
		return &RepoInfo{
			IsPrivate:        false,
			Name:             "",
			OwnerName:        "",
			OrganizationName: ev.GetOrganization().GetLogin(),
		}, nil
	case *github.MembershipEvent:
		return &RepoInfo{
			IsPrivate:        true,
			Name:             "",
			OwnerName:        "",
			OrganizationName: ev.GetOrg().GetLogin(),
		}, nil
	case *github.PushEvent:
		return &RepoInfo{
			IsPrivate:        ev.GetRepo().GetPrivate(),
			Name:             ev.GetRepo().GetName(),
			OwnerName:        ev.GetRepo().GetOwner().GetLogin(),
			OrganizationName: ev.GetOrganization().GetLogin(),
		}, nil
	default:
		if roev, ok := ev.(repoOrgEvent); ok {
			return &RepoInfo{
				IsPrivate:        roev.GetRepo().GetPrivate(),
				Name:             roev.GetRepo().GetName(),
				OwnerName:        roev.GetRepo().GetOwner().GetLogin(),
				OrganizationName: roev.GetOrganization().GetLogin(),
			}, nil
		}
		if rev, ok := ev.(repoEvent); ok {
			return &RepoInfo{
				IsPrivate:        rev.GetRepo().GetPrivate(),
				Name:             rev.GetRepo().GetName(),
				OwnerName:        rev.GetRepo().GetOwner().GetLogin(),
				OrganizationName: "",
			}, nil
		}
		return nil, nil
	}
}

// repoEvent is used to filter on GitHub Events which
// are for a repository.
type repoOrgEvent interface {
	GetRepo() *github.Repository
	GetOrganization() *github.Organization
}

type repoEvent interface {
	GetRepo() *github.Repository
}
