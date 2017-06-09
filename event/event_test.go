package event

import (
	"encoding/json"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

type EventMock struct {
	RoutingKey string        `json:"routing_key"`
	Action     string        `json:"event_action"`
	DedupKey   string        `json:"dedup_key,omitempty"`
	Images     []interface{} `json:"images,omitempty"`
	Client     string        `json:"client,omitempty"`
	ClientURL  string        `json:"client_url,omitempty"`
	Payload    *PayloadMock  `json:"payload,omitempty"`
}

type PayloadMock struct {
	Summary   string      `json:"summary"`
	Source    string      `json:"source"`
	Severity  string      `json:"severity"`
	Timestamp string      `json:"timestamp,omitempty"`
	Component string      `json:"component,omitempty"`
	Group     string      `json:"group,omitempty"`
	Class     string      `json:"class,omitempty"`
	Details   interface{} `json:"custom_details,omitempty"`
}

type ResponseMock struct {
	RoutingKey  string `json:"routing_key"`
	DedupKey    string `json:"dedup_key"`
	EventAction string `json:"event_action"`
}

func TestTriggerEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://events.pagerduty.com/v2/enqueue",
		func(req *http.Request) (*http.Response, error) {
			event := &EventMock{}
			if err := json.NewDecoder(req.Body).Decode(&event); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}

			eventResponse := &ResponseMock{
				RoutingKey:  event.RoutingKey,
				DedupKey:    event.DedupKey,
				EventAction: event.Action,
			}

			resp, err := httpmock.NewJsonResponse(202, *eventResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil

		},
	)

	testEvent := &Event{
		RoutingKey: "testServiceKey",
		Action:     "trigger",
		DedupKey:   "dedupeme",
		Payload: &Payload{
			Summary:   "Summary",
			Source:    "event_test",
			Severity:  "critical",
			Timestamp: "2006-01-02T15:04:05.000Z",
			Component: "tests",
			Group:     "Group",
			Class:     "Class",
			Details:   "myDetails",
		},
	}

	manageResponse, err := ManageEvent(*testEvent)
	if err != nil {
		t.Fail()
	}

	if manageResponse.DedupKey != testEvent.DedupKey {
		t.Fail()
	}

	if manageResponse.EventAction != testEvent.Action {
		t.Fail()
	}

	if manageResponse.RoutingKey != testEvent.RoutingKey {
		t.Fail()
	}
}
