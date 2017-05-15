package event

import (
	"encoding/json"
	"bytes"
	"net/http"
	"fmt"
)

type Event struct {
	RoutingKey string	   `json:"routing_key"`
	Action	   string	   `json:"event_action"`
	DedupKey   string	   `json:"dedup_key,omitempty"`
	Images	   []interface{}   `json:"images,omitempty"`
	Client	   string	   `json:"client,omitempty"`
	ClientURL  string          `json:"client_url,omitempty"`
	Payload	   *EventPayload   `json:"payload,omitempty"`
}

type EventPayload struct {
	Summary	   string	    `json:"summary"`
	Source	   string	    `json:"source"`
	Severity   string	    `json:"severity"`
	Timestamp  string	    `json:"timestamp,omitempty"`
	Component  string	    `json:"component,omitempty"`
	Group	   string	    `json:"group,omitempty"`
	Class	   string	    `json:"class,omitempty"`
	Details	   interface{}	    `json:"custom_details,omitempty"`
}

type EventResponse struct {
	RoutingKey 	string `json:"routing_key"`
	DedupKey	string `json:"dedup_key"`
	EventAction	string `json:"event_action"`
}

const eventEndPoint = "https://events.pagerduty.com/v2/enqueue"

func ManageEvent(e Event) (*EventResponse, error) {
	data, err :=json.Marshal(e)
	if err !=nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", eventEndPoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
	}
	var eventResponse EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventResponse); err != nil {
		return nil, err
	}
	return &eventResponse, nil
}
