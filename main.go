package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"balancer-manager/apache"
)

var version string

func main() {
	hostnameFlag := flag.String("hostname", "http://localhost", "A hostname searving the /balancer-manager URI")
	actionFlag := flag.String("action", "status", "Action disable|enable|status")
	hostsFlag := flag.String("hosts", "", "Comma-separated list of workers")
	portFlag := flag.String("port", "80", "Workers port")
	uriFlag := flag.String("uri", "", "Application uri")
	debugFlag := flag.Bool("debug", false, "Enables debugging")

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			fmt.Println(version)
			os.Exit(0)
		}
	}

	flag.Parse()

	hosts := strings.Split(*hostsFlag, ",")
	instance := apache.Apache{
		Url:   *hostnameFlag + "/balancer-manager",
		Debug: *debugFlag,
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
