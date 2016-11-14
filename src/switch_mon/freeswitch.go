package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "strconv"
  "github.com/0x19/goesl"
)

type Freeswitch struct {
  host            string
  port            uint
  password        string
  connect_timeout int
  config          ConfigJson
  client          *goesl.Client
  riemannChan     chan RiemannEvent
}

func NewFreeswitch(host string, port uint, password string, connect_timeout int, config ConfigJson, riemannChan chan RiemannEvent) *Freeswitch {
  return &Freeswitch{
    host:             host,
    port:             port,
    password:         password,
    connect_timeout:  connect_timeout,
    config:           config,
    riemannChan:      riemannChan,
  }
}

func (fs *Freeswitch) startClient() {
	client, err := goesl.NewClient(fs.host, fs.port, fs.password, fs.connect_timeout)
	if err != nil {
		log.Printf("Error while creating new client: %s\n", err)
  } else {
	  // Apparently all is good... Let us now handle connection :)
	  // We don't want this to be inside of new connection as who knows where it my lead us.
	  // Remember that this is crutial part in handling incoming messages. This is a must!
	  go client.Handle()
		fs.client = &client
    fs.setHooks()
	}
}

func (fs *Freeswitch) setHooks() {
	for _, hook := range fs.config.Hooks {
		regString := fmt.Sprintf("events json %s", hook.Event)
		fs.client.Send(regString)
	}
}

func (fs *Freeswitch) loop() {
  defer func() {
    if fs.client != nil {
      fs.client.Close()
    }
  }()
	for {
    if fs.client == nil {
      time.Sleep(1000 * time.Millisecond)
      fs.startClient()
    } else {
		  msg, err := fs.client.ReadMessage()
		  if err != nil {
		    // If it contains EOF, we really dont care...
		    log.Printf("Error while reading Freeswitch message: %v", err)
		    if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
		      log.Printf("Error while reading Freeswitch message: %s\n", err)
		    }
        fs.client = nil
      } else {
		    event := msg.Headers["Event-Name"]
		    for _, hook := range fs.config.Hooks {
		      if hook.Event == event {
		        attributes := make(map[string]string)
		        for _, attrName := range hook.Attributes {
		          attributes[attrName] = msg.Headers[attrName]
		        }

		        if len(hook.Metrics) > 0 {
		          for _, metricName := range hook.Metrics {
		            if metric, err := strconv.ParseFloat(msg.Headers[metricName], 32); err == nil {
		              fs.riemannChan <- RiemannEvent{
		                Service:    fmt.Sprintf("switchmon %s %s", hook.Service, metricName),
		                Metric:     metric,
		                Attributes: attributes,
		              }
		            }
		          }
		        } else {
		          fs.riemannChan <- RiemannEvent{
		            Service:    fmt.Sprintf("switchmon %s", hook.Service),
		            Metric:     0,
		            Attributes: attributes,
		          }
		        }
		      }
		    }
		  }
    }
	}
}

func (fs *Freeswitch) Run() {
  fs.startClient()
  fs.loop()
}

