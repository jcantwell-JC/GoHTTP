package main

import (
  "testing"
  "time"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
  "strings"
  "log"
  "bytes"
)

func runHashEndpoint() *httptest.Server{
  handler := &HashHandler{}
  ts := httptest.NewServer(handler)
  return ts
}

func runStatsEndpoint() *httptest.Server{
  handler := &StatsHandler{}
  ts := httptest.NewServer(handler)
  return ts
}

//////////////////////////////////////////////
//////////////// Unit Tests //////////////////
//////////////////////////////////////////////

func TestGenHash(t *testing.T) {
  hash := generate_hash("angryMonkey")
  if hash != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
   t.Errorf("Hash was incorrect. got %s, want ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q== ", hash)
  }
}

func TestAddSummedResponseTimeAppends1ElementCorrectly(t *testing.T) {
  // add repsonseTime
  summedHashResponseTimes = generateSummedHashResponseTimes([]string{"1000ns"})

  // assert responseTimes slice is size 1
  if len(summedHashResponseTimes) != 1 {
    t.Errorf("Expected length summedHashResponseTimes to be 1. got %d", len(summedHashResponseTimes))
  }
  // assert responseTime value is 1 microSecond
  if summedHashResponseTimes[0] != 1 {
    t.Errorf("Expected length summedHashResponseTimes to be 1. got %d", summedHashResponseTimes[0])
  }
}

func TestAddSummedResponseTimeAppends2ElementCorrectly(t *testing.T) {
  // add 2 repsonseTimes
  summedHashResponseTimes = generateSummedHashResponseTimes([]string{"1000ns", "1000ns"})

  // assert responseTimes slice is size 2
  if len(summedHashResponseTimes) != 2 {
    t.Errorf("Expected length summedHashResponseTimes to be 2. got %d", len(summedHashResponseTimes))
  }
  // assert first responseTime value is 1 microSecond
  if summedHashResponseTimes[0] != 1 {
    t.Errorf("Expected length summedHashResponseTimes to be 1. got %d", summedHashResponseTimes[0])
  }
  // assert second responseTime value is 2 microSecond; 1 + 1
  if summedHashResponseTimes[1] != 2 {
    t.Errorf("Expected length summedHashResponseTimes to be 2. got %d", summedHashResponseTimes[1])
  }
}

func TestCalcAverageResponseTimeNoTimes(t *testing.T) {
  var summedHashResponseTimes []time.Duration
  avg := calcAverageResponseTime(summedHashResponseTimes)
  if avg != 0 {
    t.Errorf("Expected calcAverageResponseTime to return 0. got %f", avg)
  }
}

func TestCalcAverageResponseTimeIsCorrectSize1(t *testing.T) {
  // initalize variables
  summedHashResponseTimes = generateSummedHashResponseTimes([]string{"1000ns"})

  // now test
  avg := calcAverageResponseTime(summedHashResponseTimes)
  if avg != 1 {
    t.Errorf("Expected calcAverageResponseTime to return 1. got %f", avg)
  }
}

func TestCalcAverageResponseTimeIsCorrectSize2(t *testing.T) {
  summedHashResponseTimes = generateSummedHashResponseTimes([]string{"1000ns", "2000ns"})

  // now test
  avg := calcAverageResponseTime(summedHashResponseTimes)
  if avg != 1.5 {
    t.Errorf("Expected calcAverageResponseTime to return 1.5. got %f", avg)
  }
}

func TestCalcAverageResponseTimeIsCorrectSizeMulti(t *testing.T) {
  summedHashResponseTimes = generateSummedHashResponseTimes([]string{"1000ns", "2000ns", "1000ns", "2000ns"})

  // now test
  avg := calcAverageResponseTime(summedHashResponseTimes)
  if avg != 1.5 {
    t.Errorf("Expected calcAverageResponseTime to return 1.5. got %f", avg)
  }
}

//////////////////////////////////////////////
/////////////// Helper Methods ///////////////
//////////////////////////////////////////////

func generateSummedHashResponseTimes(timeDurrations []string) []time.Duration{
  // initalize variables
  var summedHashResponseTimes []time.Duration
  for i := 0; i < len(timeDurrations); i++ {
    ns, _ := time.ParseDuration(timeDurrations[i])
		summedHashResponseTimes = addSummedResponseTime(ns, summedHashResponseTimes)
	}
  return summedHashResponseTimes
}


//////////////////////////////////////////////
///////////// Integration Tests //////////////
//////////////////////////////////////////////

func TestPostHashEndpointSucceeds(t *testing.T) {
  ts := runHashEndpoint()
  defer ts.Close()

  // Build the request
  resp, err :=  http.Post(ts.URL + "/hash", "application/x-www-form-urlencoded", strings.NewReader("password=angryMonkey"))
	if err != nil {
		t.Errorf("Expected no error. Error: %s", err)
	}
  if resp.StatusCode != 200 {
    t.Errorf("Expected 200 error code. Got %d", resp.StatusCode)
  }
  greeting, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
  // returns the urlencoded version of the hash. i.e + are replaced with - and / are replaced with _
  if bytes.Equal(greeting, []byte("ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==")) {
    t.Errorf("Expected POST /hash to return ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==. Got %s", greeting)
  }
}

func TestGetHashEndpointFails(t *testing.T) {
  ts := runHashEndpoint()
  defer ts.Close()
  // Build the request
	resp, err :=  http.Get(ts.URL + "/hash")
  if err != nil {
		t.Errorf("Expected no error. Error: %s", err)
	}
  if resp.StatusCode != 404 {
    t.Errorf("Expected 200 error code. Got %d", resp.StatusCode)
  }
}

func TestPostStatsEndpointFails(t *testing.T) {
  ts := runStatsEndpoint()
  defer ts.Close()
  // Build the request
	resp, err :=  http.Post(ts.URL + "/stats", "application/x-www-form-urlencoded", strings.NewReader("somestring"))
  if err != nil {
		t.Errorf("Expected no error. Error: %s", err)
	}
  if resp.StatusCode != 404 {
    t.Errorf("Expected 200 error code. Got %d", resp.StatusCode)
  }
}
