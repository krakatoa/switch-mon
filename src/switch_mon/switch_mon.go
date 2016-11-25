package main

import (
	"flag"
	"log"

  "github.com/op/go-logging"
)

func main() {
	flagConfig := flag.String("config", "./config.json", "Config file")
	flagRiemannHost := flag.String("riemann-host", "localhost", "Riemann host")
	flagRiemannPort := flag.Int("riemann-port", 5555, "Riemann port")
	flagRiemannReconnectInterval := flag.Int("riemann-reconnect-interval", 5, "Riemann reconnect interval")

  var (
    fshost   = flag.String("fshost", "localhost", "Freeswitch hostname. Default: localhost")
    fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
    password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
    timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
  )

	flag.Parse()

  logging.SetLevel(logging.WARNING, "goesl")

	riemann := NewRiemann(*flagRiemannHost, *flagRiemannPort, *flagRiemannReconnectInterval)
	go riemann.Run()

	config := ReadConfig(*flagConfig)
	for _, hook := range config.Hooks {
		log.Printf("Hook: %s\n", hook.Event)
		for _, attr := range hook.Attributes {
			log.Printf("  attr: %s\n", attr)
		}
	}

  freeswitch := NewFreeswitch(*fshost, *fsport, *password, *timeout, config, riemann.InputChan)
  go freeswitch.Run()

  select {}
}
