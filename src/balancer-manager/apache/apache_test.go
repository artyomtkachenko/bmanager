package apache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_check(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	check(errors.New("foobar"))
}

func Test_sendRequest_success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Obama")
	}))
	defer ts.Close()

	body, _ := sendRequest(ts.URL)
	if !strings.Contains(string(body), "Hello Obama") {
		t.Errorf("Expected: body to be Hello, got [%s]\n", body)
	}
}

func Test_sendRequest_fail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	_, err := sendRequest(ts.URL)
	if err == nil {
		t.Errorf("Could not execute %s\n", ts.URL)
	}
}

func Test_setDetailsFromUri(t *testing.T) {
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

func Test_getWorkerUrl(t *testing.T) {
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

func Test_parseStatusHtmlPage_ohs(t *testing.T) {
	inst := Apache{
		Url: "http://localhost/balancer-manager",
	}
	body, _ := ioutil.ReadFile("testdata/bm_ohs.html")
	(inst).parseStatusHtmlPage(strings.NewReader(string(body)))
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

func Test_parseStatusHtmlPage_vanilla(t *testing.T) {
	inst := Apache{
		Url: "http://localhost/balancer-manager",
	}
	body, _ := ioutil.ReadFile("testdata/bm_vanilla.html")
	err := (inst).parseStatusHtmlPage(strings.NewReader(string(body)))
	if (inst).kind != "vanilla" {
		t.Errorf("Expected: vanilla, got %s\n", (inst).kind)
	}
	if err != nil {
		t.Errorf("Expected: nil, got %s\n", err)
	}
}

func Test_generateReport(t *testing.T) {
	testData := map[string]string{"foo": "bar"}

	generateReport(testData)
}

func Test_getStatusForAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadFile("testdata/bm_ohs.html")
		fmt.Fprint(w, string(body))
	}))
	defer ts.Close()

	inst := Apache{
		Url: ts.URL,
	}
	(inst).getStatusForAll()

	if (inst).kind != "ohs" {
		t.Errorf("Expected: server kind to be ohs, got %s\n", (inst).kind)
	}
}

func Test_generateActionUrl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadFile("testdata/bm_ohs.html")
		fmt.Fprint(w, string(body))
	}))
	defer ts.Close()

	inst := Apache{
		Url:   ts.URL,
		Debug: true,
	}
	(inst).getStatusForAll()
	url := (inst).generateActionUrl("host000003900:8083/foo", "&dw=Disable")

	if !strings.Contains(url, "&dw=Disable") {
		t.Errorf("Expected: to get &dw=Disable in a string , got %s\n", url)
	}
}

func Test_action(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadFile("testdata/bm_ohs.html")
		fmt.Fprint(w, string(body))
	}))
	defer ts.Close()

	inst := Apache{
		Url:   ts.URL,
		Debug: true,
	}
	(inst).getStatusForAll()
	realData := (inst).action("status", []string{"host000003900", "blah"}, "8083", "/foo")

	testData := map[string]string{
		"host000003900:8083/foo": "Ok",
		"blah:8083/foo":          "NO WORKER FOUND",
	}

	for key, value := range testData {
		if realData[key] != value {
			t.Errorf("Expected: to get %s but got %s\n", value, realData[key])
		}
	}
}

func Benchmark_parseStatusHtmlPage(b *testing.B) {
	inst := Apache{
		Url: "http://localhost/balancer-manager",
	}
	body, _ := ioutil.ReadFile("testdata/bm_ohs.html")

	for i := 0; i < b.N; i++ {
		(inst).parseStatusHtmlPage(strings.NewReader(string(body)))
	}

}
