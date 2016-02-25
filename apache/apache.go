package apache

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"github.com/golang/net/html"
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
	statusAll map[string]BalancerStatus
	mainUrl   string
	kind      string
}

type Apache22 struct {
	Apache
}

func (a *Apache) New(kind string, mainUrl string) {
	a.kind = kind
	a.mainUrl = mainUrl
}

func (a Apache) getDetailsFromUri(worker string, status string) (string, BalancerStatus) {
	workerFull := a.mainUrl + worker
	var w = BalancerStatus{}
	w.status = status
	w.url = workerFull

	u, err := url.Parse(workerFull)
	if err != nil {
		panic(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)
	workerUrl := m["w"][0]
	if err != nil {
		panic(err)
	}
	u, err = url.Parse(workerUrl)
	w.uri = u.Path

	return u.Host, w
}

func (a *Apache) parseStatusHtmlPage(page io.Reader) error {
	// Do not like this impementation
	z := html.NewTokenizer(page)
	a.statusAll = make(map[string]BalancerStatus)
	var (
		hrefTag   = []byte("href")
		anchorTag = []byte("a")
		tdTag     = []byte("td")
		trTag     = []byte("tr")
		hrefFound = false
		tdCount   = 0
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
				host, res := a.getDetailsFromUri(balancer, status)
				a.statusAll[host+res.uri] = res
			}

		case html.StartTagToken:
			tag, hasAttr := z.TagName()
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
			if tdCount == 5 { //Balancer status seats at td[6]
				status = string(val)
			}
		}
	}
	return nil
}

func (a Apache) getBalancerManagerStatusPage() []byte {
	response, err := http.Get(a.mainUrl + "/balancer-manager")
	if err != nil || response.StatusCode != 200 {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func (a *Apache) GetStatusForAll() {
	body := a.getBalancerManagerStatusPage()
	if err := a.parseStatusHtmlPage(strings.NewReader(string(body))); err != nil {
		panic(err)
	}
}

func (a Apache) GetStatus(hosts []string, port string, uri string) map[string]string {
	var result = make(map[string]string)
	for _, host := range hosts {
		var hostPortUri string

		if port != "80" {
			hostPortUri = host + ":" + port + uri
		} else {
			hostPortUri = host + uri
		}
		if data, ok := a.statusAll[hostPortUri]; ok {
			result[hostPortUri] = data.status
		} else {
			result[hostPortUri] = "NO WORKER FOUND"
		}
	}
	return result
}

func (a Apache) action(action string, hosts []string, port string, uri string) map[string]string {
	status := a.GetStatus(hosts, port, uri)
	var (
		act         string
		flag        string
		hostPortUri string
		result      map[string]string
	)

	result = make(map[string]string)

	if action == "disable" { //TODO make this nicer somehow
		if a.kind == "vanilla" {
			act = "&status_D=1"
			flag = "Init Dis "
		} else { // We presume it is OHS
			act = "&dw=Disable"
			flag = "Dis "
		}
	} else if action == "enable" {
		if a.kind == "vanilla" {
			act = "&status_D=0"
			flag = "Init Ok "
		} else { // We presume it is OHS
			act = "&dw=Enable"
			flag = "Ok "
		}
	}

	for _, host := range hosts {
		if port == "80" {
			hostPortUri = host + uri
		} else {
			hostPortUri = host + ":" + port + uri
		}
		if status[hostPortUri] == "NO WORKER FOUND" {
			result[hostPortUri] = status[hostPortUri]
		} else if a.statusAll[hostPortUri].status != flag {
			url := a.statusAll[hostPortUri].url + act
			// fmt.Printf("Sending %s\n", url)
			response, err := http.Get(url)
			if err == nil && response.StatusCode == 200 {
				result[hostPortUri] = action + "d"
			} else {
				panic(err)
			}
		} else {
			result[hostPortUri] = "Is already in " + flag + " state"
		}
	}
	return result
}

func (a Apache) Disable(hosts []string, port string, uri string) {
	res := a.action("disable", hosts, port, uri)
	if out, err := json.Marshal(res); err == nil {
		fmt.Println(string(out))
	}
}

func (a Apache) Enable(hosts []string, port string, uri string) {
	res := a.action("enable", hosts, port, uri)
	if out, err := json.Marshal(res); err == nil {
		fmt.Println(string(out))
	}
}
