package main

import (
  "io/ioutil"

  "encoding/json"
)

type HookJson struct {
  Event       string    `json:"event"`
  Metrics     []string  `json:"metrics"`
  Attributes  []string  `json:"attributes"`
  Service     string    `json:"service"`
}

type ConfigJson struct {
  Hooks   []HookJson    `json:"hooks"`
}

func ReadConfig(configPath string) ConfigJson {
  dat, err := ioutil.ReadFile(configPath)
  if err != nil {
    panic(err)
  }
  // log.Printf("%s\n", string(dat))

  config := ConfigJson{}
  json.Unmarshal([]byte(dat), &config)

  // for _, hook := range config.Hooks {
  //   log.Printf("Hook: %s\n", hook.Event)
  //   for _, attr := range hook.Attributes {
  //     log.Printf("  attr: %s\n", attr)
  //   }
  // }

  // log.Printf("config: %v", config.Hooks)
  return config
}
