package hello

import (
  "fmt"
  "http"
  "json"
)

func init() {
  http.HandleFunc("/", handler)
}

type Success struct {
  Message string `json:"message"`
  Value int `json:"value"`
}

func handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-type", "application/json; charset=utf-8")
  m := Success{"My success message", 22}
  b, _ := json.MarshalIndent(m, "", "  ")
  n := Success{}
  json.Unmarshal(b, &n)
  fmt.Print(n)
  fmt.Fprint(w, string(b))
}
