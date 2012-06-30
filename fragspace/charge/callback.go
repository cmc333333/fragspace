package charge

import (
  "http"
)

func init() {
  http.HandleFunc("/charge/callback", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
      case "POST": callback(w, r)
      default: http.Error(w, "404 page not found", http.StatusNotFound)
    }
  })
}

type eventReq struct {
  Id string `json:"id"`
}

func callback(w http.ResponseWriter, r *http.Request) {
  post := new(eventReq)
  if err := (*fhttp.JsonRequest)(r).Extract(eventReq); err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
  }
  req, err := http.NewRequest("GET", "https://api.stripe.com/v1/events/" + eventReq.Id, nil)
  if err != nil {
    http.Error(w, "Could not create request: " + err.String(), http.StatusInternalServerError)
  }
  req.SetBasicAuth(
}
