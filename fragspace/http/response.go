package http

import (
  "http"
  "json"
)

type Response interface {
  WriteTo(w http.ResponseWriter)
}

type JsonResponse struct {
  Msg interface{}
}
func (r JsonResponse) WriteTo(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  bytes, err := json.MarshalIndent(r.Msg, "", "  ")
  if err != nil {
    panic(ServerError(err.String()))
  }
  w.Write(bytes)
}

type Success struct{}
func (s Success) WriteTo(w http.ResponseWriter) {
}

type error struct {
  Error string `json:"error"`
}

type UserError string
func (e UserError) WriteTo(w http.ResponseWriter) {
  w.WriteHeader(http.StatusBadRequest)
  JsonResponse{error{string(e)}}.WriteTo(w)
}

type ServerError string
func (e ServerError) WriteTo(w http.ResponseWriter) {
  w.WriteHeader(http.StatusInternalServerError)
  JsonResponse{error{string(e)}}.WriteTo(w)
}

type NotFound struct{}
func (n NotFound) WriteTo(w http.ResponseWriter) {
  http.Error(w, "404 page not found", http.StatusNotFound)
}
