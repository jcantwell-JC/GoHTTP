package main

import (
  "testing"
  "time"
  //"net/http"
)

func TestGenHash(t *testing.T) {
  hash := generate_hash("angryMonkey")
  if hash != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
   t.Errorf("Hash was incorrect. got %s, want ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q== ", hash)
  }
}

func TestAddSummedResponseTimeAppendsCorrectly(t *testing.T) {
  var summedHashResponseTimes []time.Duration
  summedHashResponseTimes = addSummedResponseTime(time.Since(time.Now()), summedHashResponseTimes)
  if len(summedHashResponseTimes) != 1 {
    t.Errorf("Expected length summedHashResponseTimes to be 1. got %d", len(summedHashResponseTimes))
  }
  summedHashResponseTimes = addSummedResponseTime(time.Since(time.Now()), summedHashResponseTimes)
  if len(summedHashResponseTimes) != 2 {
    t.Errorf("Expected length summedHashResponseTimes to be 2. got %d", len(summedHashResponseTimes))
  }
}

// func TestPostHashEndpointSucceeds(t *testing.T) {
//   // Build the request
// 	resp, err :=  http.Post("http://localhost:8080/hash", "application/json", bytes.NewBuffer([]byte("angryMonkey")))
// 	if err != nil {
// 		t.Errorf("Expected no error. Error: %s", err)
// 	}
//   if resp.StatusCode != 200 {
//     t.Errorf("Expected 200 error code. Got %d", resp.StatusCode)
//   }
// }
//
// func TestGetHashEndpointFails(t *testing.T) {
//   // Build the request
// 	resp, err :=  http.Get("http://localhost:8080/hash")
//   if err != nil {
// 		t.Errorf("Expected no error. Error: %s", err)
// 	}
//   if resp.StatusCode != 404 {
//     t.Errorf("Expected 200 error code. Got %d", resp.StatusCode)
//   }
// }
