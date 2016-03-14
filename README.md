#Apache Balancar Manager in Golang
Oracle HTTP Server ( Apache 2.2 ) Balancer manager application written in Golang.

##Build
```bash
make deps && make build
```

##Usage
```bash
bin/balancer-manager -help
```
```bash
Usage of bin/balancer-manager:

-action string 
Action disable|enable|status (default "status")

-debug
Enables debugging

-hostname string
A hostname searving the /balancer-manager URI (default "")

-hosts string
Comma-separated list of workers

-port string
Workers port (default "80")

-uri string
Application uri
```

##Examples
Output status command (default action when -action is not specified)
```bash
bin/balancer-manager -hosts host1,host2,...,hostn -port 8080 -uri /foo
```

Disable action command
```bash
bin/balancer-manager -action disable -hosts host1,host2,...,hostN -port 8080 -uri /foo
```

Enable action command
```bash
bin/balancer-manager -action enable -hosts host1,host2,...,hostN -port 8080 -uri /foo
```

Example output status 
```bash
{"host1:8080/foo":"Ok", "host2:8080/foo":"NO WORKER FOUND", ... , "hostN:8080/foo":"Dis "}
``` 
