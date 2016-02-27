package main

import (
	"flag"
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
	instance := apache.Apache{
		Url: *hostnameFlag + "/balancer-manager",
	}

	switch *actionFlag {
	case "enable":
		instance.Enable(hosts, *portFlag, *uriFlag)
	case "disable":
		instance.Disable(hosts, *portFlag, *uriFlag)
	case "status":
		instance.Status(hosts, *portFlag, *uriFlag)
	}
}
