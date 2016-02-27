package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/artyomtkachenko/bmanager/apache"
	"strings"
)

func main() {
	hostnameFlag := flag.String("hostname", "http://localhost", "A hostname searving the /balancer-manager URI")
	actionFlag := flag.String("action", "status", "Action disable|enable|status")
	hostsFlag := flag.String("hosts", "", "Comma-separated list of workers")
	portFlag := flag.String("port", "80", "Workers port")
	uriFlag := flag.String("uri", "", "Application uri")
	flag.Parse()

	hosts := strings.Split(*hostsFlag, ",")
	instance := new(apache.Apache22)

	if *targetFlag == "apache22" {
		instance.Init("vanilla", *hostnameFlag)
	} else if *targetFlag == "ohs" {
		instance.Init("ohs", *hostnameFlag)
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
