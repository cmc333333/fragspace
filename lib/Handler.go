package lib

import (
  "http"
  "appengine"
)

type Handler interface {
  Post(JsonResponse, *JsonRequest)
  Get(JsonResponse, *http.Request)
  Put(JsonResponse, *JsonRequest)
  Delete(JsonResponse, *http.Request)
  http.Handler
}

type BaseHandler struct {
  Self Handler
}

//  Error handling via exceptions
func errorHandling(w http.ResponseWriter, r *http.Request) {
  if err := recover(); err != nil {
    switch e := err.(type) {
      case UserError: http.Error(w, e.String(), http.StatusBadRequest)
      case ServerError:
        c := appengine.NewContext(r)
        c.Errorf(e.String())
        http.Error(w, "", http.StatusInternalServerError)
      default: panic(err)
    }
  }
}

func (handler *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  defer errorHandling(w, r)
  switch r.Method {
    case "POST": handler.Self.Post(JsonResponse{w}, (*JsonRequest)(r))
    case "GET": handler.Self.Get(JsonResponse{w}, r)
    case "PUT": handler.Self.Put(JsonResponse{w}, (*JsonRequest)(r))
    case "DELETE": handler.Self.Delete(JsonResponse{w}, r)
    default: http.NotFound(w, r)
  }
}
func (handler *BaseHandler) Post(w JsonResponse, r *JsonRequest) {
  http.NotFound(w, (*http.Request)(r))
}
func (handler *BaseHandler) Get(w JsonResponse, r *http.Request) {
  http.NotFound(w, r)
}
func (handler *BaseHandler) Put(w JsonResponse, r *JsonRequest) {
  http.NotFound(w, (*http.Request)(r))
}
func (handler *BaseHandler) Delete(w JsonResponse, r *http.Request) {
  http.NotFound(w, r)
}
