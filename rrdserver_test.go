package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	filename := "config.toml"
	ReadConfigFile(filename)

	if config.Server.RrdPath != "sample/" {
		t.Fatalf("Error reading config file %s", filename)
	}

	err := ReadConfigFile("fugahoge.toml")
	if err == nil {
		t.Fatalf("The file %s is invalid but error didn't happen.", err)
	}
}

func TestHello(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(hello))
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("Status code is not 200 but %d.", r.StatusCode)
		return
	}

	if "{\"message\":\"hello\"}" != string(data) {
		t.Fatalf("Data Error. %v", string(data))
	}
}

func TestSearch(t *testing.T) {
	requestJSON := `{"target":"hoge"}`
	reader := strings.NewReader(requestJSON)

	ReadConfigFile("config.toml")

	ts := httptest.NewServer(http.HandlerFunc(search))
	defer ts.Close()

	r, err := http.Post(ts.URL, "application/json", reader)
	if err != nil {
		t.Fatalf("Error at a GET request. %v", err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("Status code is not 200 but %d.", r.StatusCode)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var searchResponse []string
	decoder.Decode(&searchResponse)

	if len(searchResponse) <= 0 {
		t.Fatalf("Data Error. %v", searchResponse)
	} else {
		t.Log(searchResponse)
	}
}

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(query))
	defer ts.Close()
	client := &http.Client{}

	// Test for an OPTIONS request
	reqOptions, err := http.NewRequest("OPTIONS", ts.URL, nil)
	if err != nil {
		t.Fatalf("Error at an OPTIONS request. %v", err)
	}
	r, err := client.Do(reqOptions)

	if r.StatusCode != 200 {
		t.Fatalf("Status code is not 200 but %d.", r.StatusCode)
		return
	}

	if r.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("Header Access-Control-Allow-Origin is invalid. %s", r.Header.Get("Access-Control-Allow-Origin"))
	}
	if r.Header.Get("Access-Control-Allow-Headers") != "accept, content-type" {
		t.Fatalf("Header Access-Control-Allow-Headers is invalid. %s", r.Header.Get("Access-Control-Allow-Headers"))
	}
	if r.Header.Get("Access-Control-Allow-Methods") != "GET,POST,HEAD,OPTIONS" {
		t.Fatalf("Header Access-Control-Allow-Methods is invalid. %s", r.Header.Get("Access-Control-Allow-Methods"))
	}

	// Test for a POST request
	requestJSON := `{
	  "panelId":1,
	  "range":{
	    "from":"2017-01-17T05:14:42.237Z",
	    "to":"2017-01-18T05:14:42.237Z",
	    "raw":{
	      "from":"now-24h",
	      "to":"now"
	    }
	  },
	  "rangeRaw":{
	    "from":"now-24h",
	    "to":"now"
	  },
	  "interval":"1m",
	  "intervalMs":60000,
	  "targets":[
	    {"target":"sample:ClientJobsIdle","refId":"A","hide":false,"type":"timeserie"},
	    {"target":"sample:ClientJobsRunning","refId":"B","hide":false,"type":"timeserie"}
	  ],
	  "format":"json",
	  "maxDataPoints":1812
	}`

	reader := strings.NewReader(requestJSON)
	r, err = http.Post(ts.URL, "application/json; charset=utf-8", reader)
	if err != nil {
		t.Fatalf("Error at an POST request. %v", err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("Status code is not 200 but %d.", r.StatusCode)
		return
	}

	if r.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("Header Access-Control-Allow-Origin is invalid. %s", r.Header.Get("Access-Control-Allow-Origin"))
	}
	if r.Header.Get("Access-Control-Allow-Headers") != "accept, content-type" {
		t.Fatalf("Header Access-Control-Allow-Headers is invalid. %s", r.Header.Get("Access-Control-Allow-Headers"))
	}
	if r.Header.Get("Access-Control-Allow-Methods") != "GET,POST,HEAD,OPTIONS" {
		t.Fatalf("Header Access-Control-Allow-Methods is invalid. %s", r.Header.Get("Access-Control-Allow-Methods"))
	}

	decoder := json.NewDecoder(r.Body)
	var queryRequest QueryRequest
	err = decoder.Decode(&queryRequest)
	if err != nil {
		fmt.Println("error in query 1")
		fmt.Println(err)
	}

	if len(queryRequest.Targets) < 1 {
		t.Fatalf("Response is empty.")
	}
}

func TestAnnotations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(annotations))
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error at a GET request. %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	if r.StatusCode != 200 {
		t.Fatalf("Status code is not 200 but %d.", r.StatusCode)
		return
	}

	if "{\"message\":\"annotations\"}" != string(data) {
		t.Fatalf("Data Error. %v", string(data))
	}
}
