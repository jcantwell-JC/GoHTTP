package main

import (
  "testing"
  // "net/http"
  // "bytes"
)

func TestGenHash(t *testing.T) {
  hash := generate_hash("angryMonkey")
  if hash != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
   t.Errorf("Hash was incorrect. got %s, want ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q== ", hash)
  }
}

func TestWriteErrorMsg(t *testing.T) {

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
