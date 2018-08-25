package main

import (
    "fmt"
    "crypto/sha512"
    "encoding/base64"
    "log"
    "net/http"
    "time"
    "encoding/json"
    "context"
)

//////////////////////////////////////////////
///////////// Global Variables ///////////////
//////////////////////////////////////////////

// InProgress tracker for hashing; used in shutdown to ensure all hashing work has completed
var hashInProgress = false;
// tracker for the total number of times the /hash endpoint has been called
var totalHashCalls = 0;

// defines an error message structure- to make an error just a little more pretty
type ErrorMessage struct {
    Error string
}

// stats endpoint return struct
type Stats struct {
    Total int
    AverageTime string
}

//////////////////////////////////////////////
////////////////// Handlers //////////////////
//////////////////////////////////////////////

func hashHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
      case "POST":
        totalHashCalls++ // we've been called. increment
        r.ParseForm()
        hashInProgress = true;
        fmt.Printf("Waiting 5 sec..\n")
        time.Sleep(time.Duration(5)*time.Second) // Pause for 5 seconds
        hash := generate_hash(r.Form["password"][0])
        hashInProgress = false;
        write200Msg(w, []byte(hash))
        defer elapsed("/hash")()
      default:
        writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
  }
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      m := Stats{totalHashCalls, ""}
      jsonMessage, err := json.Marshal(m) // create json message with password hash
      if err != nil {
        writeErrorMsg(w, "Issue fetching data", http.StatusNotFound)
      } else {
        write200Msg(w, jsonMessage)
      }
    default:
      writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
  }
}

func main() {
  srv := &http.Server{Addr: ":8080"}

  http.HandleFunc("/hash", hashHandler)
  http.HandleFunc("/stats", statsHandler)
  http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
      case "GET":
          for true { // continue looping until hash is not in progress.
            if !hashInProgress {
              fmt.Printf("OK... shutting down\n")
              if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
                  log.Fatal(err)
              }
            }
          }
      default:
        writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
    }
  })

  fmt.Printf("Starting server\n")
  if err := srv.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}

//////////////////////////////////////////////
/////////////// Helper Methods ///////////////
//////////////////////////////////////////////

func write200Msg(w http.ResponseWriter, message []byte) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(message)
}

func writeErrorMsg(w http.ResponseWriter, message string, statusCode int) {
  m := ErrorMessage{message}
  jsonMessage, err := json.Marshal(m)
  if err != nil {
    jsonMessage = []byte("\"Error\": \"\"}")
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  w.Write(jsonMessage)
}

func generate_hash(s string) string {
    sha_512 := sha512.New()
    sha_512.Write([]byte(s))
    sha := base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
    fmt.Printf("returing sha_512: %s\n",sha)
    return sha+"\n"
}

func elapsed(what string) func() {
    start := time.Now()
    return func() {
        mytime := time.Since(start)
        fmt.Printf("%s took %v\n", what, mytime)
    }
}
