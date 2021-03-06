package handlers

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
    Total int
    Average float64
}

// initialize empty slice of time.Duration.
// This will track the running time total it takes for the /hash endpoint to return
// used in the /stats endpoint
var summedHashResponseTimes []time.Duration

//////////////////////////////////////////////
///////////// Handlers ///////////////
//////////////////////////////////////////////

type HashHandler struct {}
// needs a ServeHTTP method from HandlerFunc Interface
func (h *HashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
      case "POST":
        start := time.Now() // capture starting time
        r.ParseForm()
        formData := r.Form["password"]
        if formData == nil {
          writeErrorMsg(w, "Missing input data in request", 400)
          return
        }
        if len(formData) != 1 {
          writeErrorMsg(w, "Bad input data in request", 400)
          return
        }
        hashInProgress = true;
        fmt.Printf("Waiting 5 sec before returning hash\n")
        time.Sleep(time.Duration(5)*time.Second) // Pause for 5 seconds
        hash := generate_hash(formData[0])
        hashInProgress = false;
        //fmt.Printf("returning hash %s\n", hash)
        write200Msg(w, []byte(hash))
        elapsed := time.Since(start) // caculate how much time has passed
        summedHashResponseTimes = addSummedResponseTime(elapsed, summedHashResponseTimes) // add elapsed time to summedHashResponseTimes slice
      default:
        writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
  }
}

type StatsHandler struct {}
// needs a ServeHTTP method from HandlerFunc Interface
func (s *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      m := Stats{Total: len(summedHashResponseTimes), Average: calcAverageResponseTime(summedHashResponseTimes)}
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

type ShutdownHandler struct {
  Srv *http.Server // takes an httpServer
}
// needs a ServeHTTP method from HandlerFunc Interface
func (s *ShutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
        for true { // continue looping until hash is not in progress.
          if !hashInProgress {
            fmt.Printf("Received shutdown request... shutting down\n")
            if err := s.Srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
                log.Fatal(err)
            }
          }
        }
    default:
      writeErrorMsg(w, r.Method + " is not supported", http.StatusNotFound)
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
  // the /hash endpoint has never been hit, just return 0
  return float64(0)
}
