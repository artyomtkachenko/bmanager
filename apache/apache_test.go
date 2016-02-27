package apache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Obama")
	}))
	defer ts.Close()

	body, _ := sendRequest(ts.URL)
	if !strings.Contains(string(body), "Hello Obama") {
		t.Errorf("Expected: body to be Hello, got [%s]\n", body)
	}
}

func TestGetDetailsFromUri(t *testing.T) {
	inst := Apache{
		Url: "http://localhost/balancer-manager",
	}
	host, worker := (inst).getDetailsFromUri("/balancer-manager/?b=foo&w=http://host1:8083/foo&nonce=2afae676", "Ok")
	expectUrl := "http://localhost/balancer-manager/balancer-manager/?b=foo&w=http://host1:8083/foo&nonce=2afae676"
	if worker.url != expectUrl {
		t.Errorf("Expected: %s, got %s\n", expectUrl, worker.url)
	}
	if host != "host1:8083" {
		t.Errorf("Expected: host1:8083, got %s\n", host)
	}
}

func TestGetWorkerUrl(t *testing.T) {
	inst := Apache{
		Url: "http://localhost/balancer-manager",
	}
	testData := []struct {
		h        string
		p        string
		u        string
		realData string
	}{
		{"host1", "80", "/foo", "host1/foo"},
		{"host1", "8080", "/foo", "host1:8080/foo"},
	}
	for _, td := range testData {
		if res := (inst).getWorkerUrl(td.h, td.p, td.u); res != td.realData {
			t.Errorf("realData: %s, got %s\n", td.realData, td.h+td.p+td.u)
		}

	}
}

func TestParseStatusHtmlPage(t *testing.T) {
	inst := Apache{
		Url: "http://localhost/balancer-manager",
	}
	body, _ := ioutil.ReadFile("testdata/bm.html")
	(inst).parseStatusHtmlPage(strings.NewReader(string(body)))
	/* fmt.Printf("%+v\n", inst) */
	testData := map[string]BalancerStatus{
		"host000003900:8083/foo": BalancerStatus{
			url:    "http://localhost/balancer-manager/balancer-manager/?b=foo&w=http://host000003900:8083/foo&nonce=2afae676-da25-11e5-bf61-cff14dcff070",
			uri:    "/foo",
			status: "Ok",
		},
		"host000004000:8083/foo": BalancerStatus{
			url:    "http://localhost/balancer-manager/balancer-manager/?b=foo&w=http://host000004000:8083/foo&nonce=2afae676-da25-11e5-bf61-cff14dcff070",
			uri:    "/foo",
			status: "Dis",
		},
		"host000013700:8084/bar": BalancerStatus{
			url:    "http://localhost/balancer-manager/balancer-manager/?b=bar&w=http://host000013700:8084/bar&nonce=2afae676-da25-11e5-bf61-cff14dcff070",
			uri:    "/bar",
			status: "Ok",
		},
		"host000013800:8084/bar": BalancerStatus{
			url:    "http://localhost/balancer-manager/balancer-manager/?b=bar&w=http://host000013800:8084/bar&nonce=2afae676-da25-11e5-bf61-cff14dcff070",
			uri:    "/bar",
			status: "Okay",
		},
	}

	realData := (inst).status
	fmt.Printf("%+v\n", realData)
	for key, value := range testData {
		if realData[key].uri != value.uri {
			t.Errorf("Expected: %s, got %s\n", value.uri, realData[key].uri)
		}
		if realData[key].url != value.url {
			t.Errorf("Expected: %s, got %s\n", value.url, realData[key].url)
		}
		if realData[key].status != value.status {
			t.Errorf("Expected: %s, got %s\n", value.status, realData[key].status)
		}
	}
}

// func TestGetStatus(t *testing.T) {
// 	inst := Apache{
// 		Url: "http://localhost/balancer-manager",
// 	}
// 	body, _ := ioutil.ReadFile("testdata/bm.html")
// 	(inst).parseStatusHtmlPage(strings.NewReader(string(body)))
// 	hosts := []string{"host000013800", "obama"}
// 	res := inst.GetStatus(hosts, "8084", "/bar")
// 	if res["host000013800:8084/bar"] != "Okay" {
// 		t.Errorf("Expected Okay, got: %s\n", res["host000013800:8084/bar"])
// 	}
// 	if res["obama:8084/bar"] != "NO WORKER FOUND" {
// 		t.Errorf("Expected NO WORKER FOUND, got: %s\n", res["obama"])
// 	}
// }
