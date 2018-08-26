package main

import (
  "testing"
  "net/http"
  "strings"
)

func TestApp(t *testing.T) {
  a := App{}
  go a.Start(":8080") // start application server on port 8080
  _, err := http.Get("http://localhost:8080/stats")
	if err != nil {
		t.Errorf("Did not expect an error but got one. err %v", err)
	}
  _, err1 := http.Post("http://localhost:8080/hash", "application/x-www-form-urlencoded", strings.NewReader("somestring"))
	if err1 != nil {
		t.Errorf("Did not expect an error but got one. err %v", err1)
	}

}
