package main

import (
  "flag"
  "log"
  "github.com/0x19/goesl"
  "strings"
  "fmt"
)

var (
  fshost   = flag.String("fshost", "localhost", "Freeswitch hostname. Default: localhost")
  fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
  password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
  timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

func main() {
  // flagRiemannHost         := flag.String("riemann-host", "localhost", "Riemann host")
  // flagRiemannPort         := flag.Int("riemann-port", 5555, "Riemann port")
  // flagRiemannSendInterval := flag.Int("riemann-interval", 5, "Riemann send interval")

  // flag.Parse()

  // riemann := NewRiemann(*flagRiemannHost, *flagRiemannPort, *flagRiemannSendInterval, store, *flagDingQueueSizeWarn, *flagDingQueueSizeCrit)
  // riemann.Run()

  config := ReadConfig()
  for _, hook := range config.Hooks {
    log.Printf("Hook: %s\n", hook.Event)
    for _, attr := range hook.Attributes {
      log.Printf("  attr: %s\n", attr)
    }
  }

  client, err := goesl.NewClient(*fshost, *fsport, *password, *timeout)
  if err != nil {
    log.Printf("Error while creating new client: %s\n", err)
    return
  }

// Apparently all is good... Let us now handle connection :)
    // We don't want this to be inside of new connection as who knows where it my lead us.
    // Remember that this is crutial part in handling incoming messages. This is a must!
    go client.Handle()

    // register all hooks
    for _, hook := range config.Hooks {
      regString := fmt.Sprintf("events json %s", hook.Event)
      client.Send(regString)
    }
    //

    for {
      msg, err := client.ReadMessage()
      if err != nil {
        // If it contains EOF, we really dont care...
        if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
          // Error("Error while reading Freeswitch message: %s", err)
          log.Printf("Error while reading Freeswitch message: %s\n", err)
        }
        break
      }

      event := msg.Headers["Event-Name"]
      // var attributes []string
      vsf := make([]string, 0)
      for _, hook := range config.Hooks {
        if hook.Event == event {
          // attributes = hook.Attributes
          for _, attr := range hook.Attributes {
            vsf = append(vsf, msg.Headers[attr])
          }
        }
      }
      log.Printf("Event: %s | Attrs: %v", event, vsf)
      // log.Printf("Got new message: %s\n", msg)

      allAttrs := make([]string, 0)
      for attrName, _ := range msg.Headers {
        allAttrs = append(allAttrs, attrName)
      }
      log.Printf("  all attrs: %v", allAttrs)
    }
}
