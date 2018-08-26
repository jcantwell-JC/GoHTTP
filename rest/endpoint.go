package main

import (
    "fmt"
    "log"
    "net/http"
    "context"
    "github.com/rdibari84/GoHTTP/handlers"
)

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
  hash := handlers.HashHandler{}
  stats := handlers.StatsHandler{}
  shutdown := handlers.ShutdownHandler{Srv: srv}
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
//////////////////// Main ////////////////////
//////////////////////////////////////////////

func main() {
  a := App{}
  a.Start(":8080") // start application server on port 8080
}
