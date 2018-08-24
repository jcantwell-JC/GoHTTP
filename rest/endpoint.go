package main

import (
    "fmt"
    "crypto/sha512"
    "encoding/base64"
    "log"
    "net/http"
    "time"
    "encoding/json"
)

type Message struct {
    Body string
    Error string
}

func generate_hash(s string) string {
    sha_512 := sha512.New()
    sha_512.Write([]byte(s))
    sha := base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
    fmt.Printf("returing sha_512: %s\n",sha)
    return sha
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
      case "POST":
        r.ParseForm()
        fmt.Printf("Waiting 5 sec..\n")
        time.Sleep(time.Duration(5)*time.Second) // Pause for 5 seconds
        m := Message{generate_hash(r.Form["password"][0]), ""}
        jsonMessage, err := json.Marshal(m)
        if err != nil {
          jsonMessage = []byte("{\"Body\":\"\",\"Error\": \"Issue fetching data\"}")
          w.Header().Set("Content-Type", "application/json")
          w.WriteHeader(http.StatusInternalServerError)
          w.Write(jsonMessage)
        } else {
          w.Header().Set("Content-Type", "application/json")
          w.WriteHeader(http.StatusOK)
          w.Write(jsonMessage)
        }

      default:
        m := Message{"", r.Method + " is not supported"}
        jsonMessage, err := json.Marshal(m)
        if err != nil {
          jsonMessage = []byte("{\"Body\":\"404 Not Found\",\"Error\": \"\"}")
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write(jsonMessage)
  }
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
    http.HandleFunc("/hash", hashHandler)
    fmt.Printf("Starting server for HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
