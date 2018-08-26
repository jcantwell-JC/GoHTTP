package main

import (
  "testing"
  "time"
  "net/http"
  "net/http/httptest"
  "io/ioutil"
  "strings"
  "bytes"
  "encoding/json"
  "math"
  "fmt"
)

/**
 FIXME: The Stats Unit Tests need to run first.
 The calls that the other tests make cause these tests to fail
*/
//////////////////////////////////////////////
////////// Stats Unit Tests //////////////////
//////////////////////////////////////////////
func TestPostStatsEndpointFails(t *testing.T) {
  ts := runStatsEndpoint()
  defer ts.Close()
  // Build the request
	resp, err :=  http.Post(ts.URL + "/stats", "application/x-www-form-urlencoded", strings.NewReader("somestring"))
  if err != nil {
		t.Errorf("Expected no error. Error: %s", err)
	}
  if resp.StatusCode != 404 {
    t.Errorf("Expected 404 error code. Got %d", resp.StatusCode)
  }
}

func TestGetStatsEndpointSucceeds(t *testing.T) {
  ts := runStatsEndpoint()
  defer ts.Close()

  ch := make(chan []byte)
  go MakeStatsRequest(t, ts, ch)

  // wait for channel to return
  for true {
    // save off the channel
    value := <-ch
    if value != nil {
      stats := Stats{}
      json.Unmarshal(value, &stats)
      if stats.NumberCalls != 0 {
        t.Errorf("Expected stat.NumberCalls to be 0. got: %d", stats.NumberCalls)
      }
      if stats.AverageTime != 0 {
        t.Errorf("Expected stat.AverageTime to be 0. got: %f", stats.AverageTime)
      }
      t.Logf("{ stats.NumberCalls: %d, stats.AverageTime: %f }", stats.NumberCalls, stats.AverageTime)
    }
    break;
  }
}
func TestGetStatsEndpointSucceedsOneHashCall(t *testing.T) {
  // start Hash Endpoint
  ts := runHashEndpoint()
  defer ts.Close()

  // make hash call
  ch := make(chan []byte)
  go MakeHashRequest(t, ts, ch)

  // wait until channel finishes
  fmt.Printf("waiting for HashRequest to return \n")
  for true {
    if (<-ch != nil) {
      // start stats endpoint
      tsStats := runStatsEndpoint()
      defer tsStats.Close()

      // Make stats request
      ch1 := make(chan []byte)
      go MakeStatsRequest(t, tsStats, ch1)

      fmt.Printf("waiting for statsRequest to return \n")
      // wait until channel finishes
      for true {
        // save off the channel
        value := <-ch1
        if value != nil {
          stats := Stats{}
          json.Unmarshal(value, &stats)
          fmt.Printf("{ stats.NumberCalls: %d, stats.AverageTime: %f }\n", stats.NumberCalls, math.Floor(stats.AverageTime))
          if stats.NumberCalls != 1 {
            t.Errorf("Expected stat.NumberCalls to be 1. got: %v", stats.NumberCalls)
          }
          if math.Floor(stats.AverageTime/1000000) != 5 {
            t.Errorf("Expected stat.AverageTime to be 5. got: %v", math.Floor(stats.AverageTime/1000000))
          }
        }
        break;
      }
    }
    break;
  }
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
  summedHashResponseTimes = generateSummedHashResponseTimes([]string{"1000ns", "2000ns", "3000ns", "4000ns"})

  if summedHashResponseTimes[3] != 10 {
    t.Errorf("Expected last summedHashResponseTimes value to be 10. got %d", summedHashResponseTimes[3])
  }
  // now test
  avg := calcAverageResponseTime(summedHashResponseTimes)
  if avg != 2.5 {
    t.Errorf("Expected calcAverageResponseTime to return 2.5. got %f", avg)
  }
}

//////////////////////////////////////////////
//////// Hash Endpoint Unit Tests ////////////
//////////////////////////////////////////////

func TestPostHashEndpointSucceeds(t *testing.T) {
  ts := runHashEndpoint()
  defer ts.Close()

  ch := make(chan []byte)
  go MakeHashRequest(t, ts, ch)

  // returns the urlencoded version of the hash. i.e + are replaced with - and / are replaced with _
  if bytes.Equal(<-ch, []byte("ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==")) {
    t.Errorf("Expected POST /hash to return ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==. Got %s", <-ch)
  }
}

func TestMultiplePostsHashEndpointSucceeds(t *testing.T) {
  ts := runHashEndpoint()
  defer ts.Close()

  // Make 10 requests
  ch := make(chan []byte)
  for i := 0; i < 10; i++ {
    go MakeHashRequest(t, ts, ch)
  }

  // check all 10 requests
  for i := 0; i < 10; i++ {
    // returns the urlencoded version of the hash. i.e + are replaced with - and / are replaced with _
    if bytes.Equal(<-ch, []byte("ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==")) {
      t.Errorf("Expected POST /hash to return ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==. Got %s", <-ch)
    }
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

//////////////////////////////////////////////
////////// Hash Shutdown Unit Tests /////////////
//////////////////////////////////////////////

func TestPostShutdownEndpointFails(t *testing.T) {
  ts := runShutdownEndpoint()
  defer ts.Close()
  // Build the request
	resp, err :=  http.Post(ts.URL + "/shutdown", "application/x-www-form-urlencoded", strings.NewReader("somestring"))
  if err != nil {
		t.Errorf("Expected no error. Error: %s", err)
	}
  if resp.StatusCode != 404 {
    t.Errorf("Expected 404 error code. Got %d", resp.StatusCode)
  }
}

func TestGetShutdownEndpointSucceeds(t *testing.T) {
  ts := runShutdownEndpoint()

  // expecting an error
  MakeShutdownRequest(t, ts)

  // call it again. expecting an error
  MakeShutdownRequest(t, ts)
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

func runHashEndpoint() *httptest.Server {
  handler := &HashHandler{}
  ts := httptest.NewServer(handler)
  return ts
}

func runStatsEndpoint() *httptest.Server {
  statshandler := &StatsHandler{}
  ts := httptest.NewServer(statshandler)
  return ts
}

func runShutdownEndpoint() *httptest.Server {
  shutdownhandlder := &ShutdownHandler{}
  ts := httptest.NewServer(shutdownhandlder)
  return ts
}


func MakeHashRequest(t *testing.T, ts *httptest.Server, ch chan<-[]byte) {
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
    t.Errorf("Expected no error. Error: %s", err)
  }

  // put the response back on the channel so the tests can check response
  ch <- greeting
}

func MakeStatsRequest(t *testing.T, ts *httptest.Server, ch chan<-[]byte) {
  // Build the request
	resp, err :=  http.Get(ts.URL + "/stats")
  if err != nil {
		t.Errorf("Expected no error. Error: %s", err)
	}
  // check error code
  if resp.StatusCode != 200 {
    t.Errorf("Expected 200 error code. Got %d", resp.StatusCode)
  }

  greeting, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()

  if err != nil {
    t.Errorf("Expected no error. Error: %s", err)
  }

  // put the response back on the channel so the tests can check response
  ch <- greeting
}

func MakeShutdownRequest(t *testing.T, ts *httptest.Server) {
  // Build the request
	_, err :=  http.Get(ts.URL + "/shutdown")
  if err == nil {
		t.Errorf("Expected an error since this endpoint shuts down the server.")
	}
}
