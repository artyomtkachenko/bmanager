package apache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type BalancerStatus struct {
	url    string
	uri    string
	status string
}

type Apache struct {
	Url        string
	status     map[string]BalancerStatus
	kind       string
	disableArg string
	enableArg  string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sendRequest(url string) ([]byte, error) {
	response, err := http.Get(url) //fires the HTTP request
	if err != nil || response.StatusCode != 200 {
		return []byte(""), errors.New("Could not execute " + url)
	}
	body, err := ioutil.ReadAll(response.Body)
	check(err)
	return body, nil
}

func (self Apache) getDetailsFromUri(worker string, status string) (string, BalancerStatus) {
	workerFull := self.Url + worker
	var w = BalancerStatus{}
	w.status = status
	w.url = workerFull

	u, err := url.Parse(workerFull)
	check(err)

	m, _ := url.ParseQuery(u.RawQuery)
	workerUrl := m["w"][0]
	check(err)

	u, err = url.Parse(workerUrl)
	check(err)
	w.uri = u.Path

	return u.Host, w
}

func (self Apache) getWorkerUrl(host string, port string, uri string) string {

	if port == "80" {
		return host + uri
	} else {
		return host + ":" + port + uri
	}
}

func (self *Apache) parseStatusHtmlPage(page io.Reader) error {
	z := html.NewTokenizer(page)
	self.status = make(map[string]BalancerStatus)
	var (
		hrefTag   = []byte("href")
		anchorTag = []byte("a")
		tdTag     = []byte("td")
		trTag     = []byte("tr")
		dtTag     = []byte("dt")
		hrefFound = false
		tdCount   = 0
		dtCount   = 0
		dtFound   = false
		balancer  string
		status    string
	)
	for {
		switch z.Next() {

		case html.ErrorToken:
			if z.Err() == io.EOF {
				return nil
			} else {
				return z.Err()
			}

		case html.EndTagToken:
			tag, _ := z.TagName()
			if bytes.Equal(trTag, tag) && hrefFound {
				tdCount = 0
				hrefFound = false
				host, res := self.getDetailsFromUri(balancer, status)
				self.status[host+res.uri] = res
			}

		case html.StartTagToken:
			tag, hasAttr := z.TagName()
			if bytes.Equal(dtTag, tag) {
				dtCount += 1
			}
			if hrefFound {
				if bytes.Equal(tdTag, tag) {
					tdCount += 1
				}
			}
			if hasAttr && bytes.Equal(anchorTag, tag) {
				key, val, _ := z.TagAttr()
				if bytes.Equal(hrefTag, key) {
					balancer = string(val)
					hrefFound = true
				}
			}

		case html.TextToken:
			val := z.Text()
			if dtCount == 1 && !dtFound {
				dtFound = true
				if strings.Contains(string(val), "Oracle-HTTP-Server") {
					self.kind = "ohs"
					self.disableArg = "&dw=Disable"
					self.enableArg = "&dw=Enable"
				} else {
					self.kind = "vanilla"
					self.disableArg = "&status_D=1"
					self.enableArg = "&status_D=0"
				}
			}
			if tdCount == 5 { //Balancer status seats at td[6]
				status = string(val)
			}
		}
	}
	return nil
}

func generateReport(result map[string]string) {
	out, err := json.Marshal(result)
	check(err)
	fmt.Println(string(out))
}

// Gets statuses for all workers
func (self *Apache) getStatusForAll() {
	body, _ := sendRequest(self.Url)
	err := self.parseStatusHtmlPage(strings.NewReader(string(body)))
	check(err)
}

//Performs enable,  disable or status actions
func (self Apache) action(action string, hosts []string, port string, uri string) {
	result := make(map[string]string)
	self.getStatusForAll()

	for _, host := range hosts {
		hostPortUri := self.getWorkerUrl(host, port, uri)
		if action == "status" {
			if data, ok := self.status[hostPortUri]; ok {
				result[hostPortUri] = data.status
			} else {
				result[hostPortUri] = "NO WORKER FOUND"
			}
		} else {
			url := self.status[hostPortUri].url + action
			_, err := sendRequest(url)
			check(err)
		}
	}
	if action == "status" {
		generateReport(result)
	}
}

//Returns the current status
func (self Apache) Status(hosts []string, port string, uri string) {
	self.action("status", hosts, port, uri)
}

//Disables workers
func (self Apache) Disable(hosts []string, port string, uri string) {
	self.action(self.disableArg, hosts, port, uri)
	self.action("status", hosts, port, uri)
}

//Enables workers
func (self Apache) Enable(hosts []string, port string, uri string) {
	self.action(self.enableArg, hosts, port, uri)
	self.action("status", hosts, port, uri)
}
