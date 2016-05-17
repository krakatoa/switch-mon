//package main
//
//import (
//  "fmt"
//  "log"
//  "strconv"
//  "time"
//
//  "github.com/amir/raidman"
//)
//
//type Riemann struct {
//  host          string
//  port          int
//  sendInterval  int
//  client        *raidman.Client
//}
//
//func (r *Riemann) startClient() {
//  clientString := fmt.Sprintf("%s:%d", r.host, r.port)
//
//  client, err := raidman.Dial("tcp", clientString)
//  if err != nil {
//    log.Println("error starting Riemann client")
//  } else {
//    r.client = client
//  }
//}
//
//func NewRiemann(riemannHost string, riemannPort int, sendInterval int) *Riemann {
//  riemann := &Riemann{
//    host:          riemannHost,
//    port:          riemannPort,
//    sendInterval:  sendInterval,
//  }
//  riemann.startClient()
//
//  return riemann
//}
//
//func (r *Riemann) send(size int) {
//  var event = &raidman.Event{
//    State:   nil,
//    Service: "service",
//    Metric:  size,
//    Ttl:     60,
//  }
//
//  err := r.client.Send(event)
//  if err != nil {
//    log.Println("error sending Riemann event: %s", err)
//    r.client.Close()
//    r.client = nil
//  }
//}
//
//func (r *Riemann) Run() {
//  ro := gorocksdb.NewDefaultReadOptions()
//
//  defer func() {
//    if r.client != nil {
//      r.client.Close()
//    }
//  }()
//  for {
//    if r.client != nil {
//      value, err := r.store.db.Get(ro, []byte("count"))
//      if err != nil {
//        log.Println("error reading count statistics")
//      }
//      count, _ := strconv.Atoi(string(value.Data()))
//      value.Free()
//
//      r.send(count)
//    } else {
//      r.startClient()
//    }
//    time.Sleep(time.Duration(r.sendInterval) * time.Second)
//  }
//}
