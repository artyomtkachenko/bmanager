package main

import (
  "fmt"
  "strings"
  "flag"
  "github.com/artyomtkachenko/apache"
  "encoding/json"
)

func main() {
  actionFlag := flag.String("action", "status", "Action disable|enable|status")
  targetFlag := flag.String("target", "apache22", "Target apache22|ohs")
  hostsFlag  := flag.String("hosts", "", "Comma-separated list of workers")
  portFlag   := flag.String("port", "80", "Workers port")
  uriFlag    := flag.String("uri", "", "Application uri")
  flag.Parse()

  hosts    := strings.Split(*hostsFlag, ",")
  instance := new(apache.Apache22)

  if *targetFlag == "apache22" {
   instance.Init("vanilla", "http://192.168.122.133")
  } else if *targetFlag == "ohs" {
    instance.Init("ohs", "http://192.168.122.133")
  }

  instance.GetStatusForAll()

  if *actionFlag == "enable" {
    instance.Enable(hosts, *portFlag, *uriFlag)
  } else if *actionFlag == "disable" {
    instance.Disable(hosts, *portFlag, *uriFlag)
  } else {
    res := instance.GetStatus(hosts, *portFlag, *uriFlag)
    if out, err := json.Marshal(res); err == nil {
      fmt.Println(string(out))
    }
  }
}
