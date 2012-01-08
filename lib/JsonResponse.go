package lib

import (
  "http"
  "json"
)

type JsonResponse struct {
  http.ResponseWriter
}

func (w JsonResponse) Success(value interface{}) {
  w.Header().Set("Content-type", "application/json; charset=utf-8")
  bytes, err := json.MarshalIndent(value, "", "  ")
  if err != nil {
    panic(ServerError{err.String()})
  }
  w.Write(bytes)
}
