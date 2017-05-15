package main

import (
	"flag"
	"log"
	"strings"
	"./event"
	"time"
)

type PRTGEvent struct {
	Probe       string
	Device      string
	Name        string
	Status      string
	Date        string
	Link        string
	Message     string
	ServiceKey  string
	IncidentKey string
	Severity    string
}

func main() {
	var probe = flag.String("probe", "local", "The PRTG probe name")
	var device = flag.String("device", "device", "The PRTG device name")
	var name = flag.String("name", "name", "The PRTG sensor name for the device")
	var status = flag.String("status", "status", "The current status for the event")
	var date = flag.String("date", "date", "The date time for the triggered event")
	var link = flag.String("linkdevice", "http://localhost", "The link to the triggering sensor")
	var message = flag.String("message", "message", "The PRTG message for the alert")
	var serviceKey = flag.String("servicekey", "myServiceKey", "The PagerDuty v2 service integration key")
	var severity = flag.String("severity", "error", "The severity level of the incident (critical, error, warning, or info)")
	flag.Parse()

	pd := &PRTGEvent{
		Probe:       *probe,
		Device:      *device,
		Name:        *name,
		Status:      *status,
		Date:        *date,
		Link:        *link,
		Message:     *message,
		ServiceKey:  *serviceKey,
		IncidentKey: *probe + "-" + *device + "-" + *name,
		Severity:    *severity,
	}

	if strings.Contains(pd.Status, "Up") || strings.Contains(pd.Status, "ended") {
		resolveEvent(pd)
	} else {
		event, err := triggerEvent(pd)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(event)
	}
}


func triggerEvent(prtg *PRTGEvent) (*event.EventResponse, error) {
	const layout = "2006-01-02T15:04:05.000Z"
	t,err := time.Parse(layout, prtg.Date)
	if err != nil {
		t = time.Now()
	}
	newEvent := &event.Event{
		RoutingKey: prtg.ServiceKey,
		Action: "trigger",
		DedupKey: prtg.IncidentKey,
		Client: "PRTG: " + prtg.IncidentKey,
		ClientURL: prtg.Link,
		Payload: &event.EventPayload{
			Summary: prtg.IncidentKey,
			Timestamp: t.Format(layout),
			Source: prtg.Link,
			Severity: prtg.Severity,
			Component: prtg.Device,
			Group: prtg.Probe,
			Class: prtg.Name,
			Details: "Link: " + prtg.Link +
				"\nIncidentKey: " + prtg.IncidentKey +
				"\nStatus: " + prtg.Status +
				"\nDate: " + prtg.Date +
				"\nMessage: " + prtg.Message,
		},
	}
	res, err := event.ManageEvent(*newEvent)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func resolveEvent(prtg *PRTGEvent) (*event.EventResponse, error) {
	triggeredEvent := &event.Event{
		RoutingKey: prtg.ServiceKey,
		Action: "resolve",
		DedupKey: prtg.IncidentKey,
	}
	res, err := event.ManageEvent(*triggeredEvent)
	if err != nil {
		return nil, err
	}
	return res, nil
}
