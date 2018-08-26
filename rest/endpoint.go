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

// defines an error message structure- to make an error a little pretty
type ErrorMessage struct {
    Error string
}

// stats endpoint return message format
type Stats struct {
    NumberCalls int
    AverageTime float64
}

// initialize empty slice of time.Duration.
// This will track the running time total it takes for the /hash endpoint to return
var summedHashResponseTimes []time.Duration

//////////////////////////////////////////////
/////// Define Http Server Structure /////////
//////////////////////////////////////////////

type App struct {
  srv http.Server
}

// start http server
func (a *App) Start(addr string) {
  srv := &http.Server{Addr: addr}

  // Create the handlers
  hash := HashHandler{}
  stats := StatsHandler{}
  shutdown := ShutdownHandler{srv}
  // now serve the handler functions
  http.HandleFunc("/hash", hash.ServeHTTP)
  http.HandleFunc("/stats", stats.ServeHTTP)
  http.HandleFunc("/shutdown", shutdown.ServeHTTP)

  fmt.Printf("Starting server\n")
  if err := srv.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}

// shutdown http server
func (a *App) Shutdown() {
  fmt.Printf("OK... shutting down\n")
  if err := a.srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
      log.Fatal(err)
  }
}


//////////////////////////////////////////////
////////////////// Handlers //////////////////
//////////////////////////////////////////////

// makes this unitestable
type HashHandler struct {}
// needs a ServeHTTP method from HandlerFunc Interface
func (h *HashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
      case "POST":
        start := time.Now() // capture starting time
        r.ParseForm()
        hashInProgress = true;
        fmt.Printf("Waiting 5 sec before returning hash\n")
        time.Sleep(time.Duration(5)*time.Second) // Pause for 5 seconds
        hash := generate_hash(r.Form["password"][0])
        hashInProgress = false;
        write200Msg(w, []byte(hash))
        elapsed := time.Since(start) // caculate how much time has passed
        summedHashResponseTimes = addSummedResponseTime(elapsed, summedHashResponseTimes) // add elapsed time to slice
      default:
        writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
  }
}

// makes this unitestable
type StatsHandler struct {}
// needs a ServeHTTP method from HandlerFunc Interface
func (s *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      m := Stats{len(summedHashResponseTimes), calcAverageResponseTime(summedHashResponseTimes)}
      fmt.Printf("returning { NumberCalls: %d, AverageTime: %f}\n", len(summedHashResponseTimes), calcAverageResponseTime(summedHashResponseTimes))
      jsonMessage, err := json.Marshal(m) // create json message with password hash
      if err != nil {
        writeErrorMsg(w, "Issue fetching data", http.StatusInternalServerError)
      } else {
        write200Msg(w, jsonMessage)
      }
    default:
      writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
  }
}

// makes this unitestable
type ShutdownHandler struct {
  srv *http.Server // takes an httpServer
}
// needs a ServeHTTP method from HandlerFunc Interface
func (s *ShutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
        for true { // continue looping until hash is not in progress.
          if !hashInProgress {
            fmt.Printf("OK... shutting down\n")
            if err := s.srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
                log.Fatal(err)
            }
          }
        }
    default:
      writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
  }
}

//////////////////////////////////////////////
//////////////////// Main ////////////////////
//////////////////////////////////////////////

func main() {
  a := App{}
  a.Start(":8080") // start application server on port 8080
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
    sha := base64.StdEncoding.EncodeToString(sha_512.Sum(nil)) // use standard endcoding instead of urlencoding. uses + and /
    return sha
}

func addSummedResponseTime(newValue time.Duration, summedHashResponseTimes []time.Duration) []time.Duration{
    microSecValue := newValue / 1000 // time.Duration stores values in nanoseconds. convert to mircoseconds
    countResponseTimes := len(summedHashResponseTimes)
    if countResponseTimes >= 1 {
      // get last summed value and caculate new summed total
      oldValue := summedHashResponseTimes[countResponseTimes-1]
      sum := oldValue + microSecValue
      summedHashResponseTimes = append(summedHashResponseTimes, sum)
    } else { // if empty, initalize with microSecond
      summedHashResponseTimes = append(summedHashResponseTimes, microSecValue)
    }
    return summedHashResponseTimes
}

func calcAverageResponseTime(summedHashResponseTimes []time.Duration) float64 {
  countResponseTimes := len(summedHashResponseTimes)
  if countResponseTimes > 0 {
    sum := summedHashResponseTimes[countResponseTimes-1] // last value is the sum
    avg := (float64(sum))/float64(countResponseTimes)
    return avg
  }
  return float64(0)
}
