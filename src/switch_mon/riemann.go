package main

import (
	"fmt"
	"log"

	"github.com/amir/raidman"
)

type Riemann struct {
	InputChan         chan RiemannEvent
	host              string
	port              int
	reconnectInterval int
	client            *raidman.Client
}

type RiemannEvent struct {
	Service    string
	Metric     interface{}
	Attributes map[string]string
}

func NewRiemann(riemannHost string, riemannPort int, reconnectInterval int) *Riemann {
	riemann := &Riemann{
		host:              riemannHost,
		port:              riemannPort,
		reconnectInterval: reconnectInterval,
		InputChan:         make(chan RiemannEvent),
	}

	return riemann
}

func (r *Riemann) startClient() {
	clientString := fmt.Sprintf("%s:%d", r.host, r.port)

	client, err := raidman.Dial("tcp", clientString)
	if err != nil {
		log.Println("error starting Riemann client")
	} else {
		r.client = client
	}
}

// switchmon channel_answer variable_rtp_audio_in_quality_percentage
func (r *Riemann) send(baseEvent RiemannEvent) {
  if r.client == nil { r.startClient() }

	var event = &raidman.Event{
		State:      "ok",
		Service:    baseEvent.Service,
		Metric:     baseEvent.Metric,
		Attributes: baseEvent.Attributes,
		Ttl:        60,
	}

	// log.Printf("Event: %v", event)
	err := r.client.Send(event)
	if err != nil {
		log.Println("error sending Riemann event: %s", err)
		r.client.Close()
		r.client = nil
	}
}

func (riemann *Riemann) Run() {
	riemann.startClient()

	defer func() {
		if riemann.client != nil {
			riemann.client.Close()
		}
	}()
	for {
		select {
		case msg := <-riemann.InputChan:
			riemann.send(msg)
		}
	}
}
