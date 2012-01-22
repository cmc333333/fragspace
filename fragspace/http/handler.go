package http

import (
  "fmt"
  "http"
  "appengine"
)

type Handler interface {
  Post(*JsonRequest) Response
  Get(*http.Request) Response
  Put(*JsonRequest) Response
  Delete(*http.Request) Response
  http.Handler
}

type BaseHandler struct {
  Self Handler
}

//  Error handling if exception occurred
func errorHandling(w http.ResponseWriter, r *http.Request) {
  if err := recover(); err != nil {
    c := appengine.NewContext(r)
    c.Errorf("%v",err)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    ServerError(fmt.Sprintf("%v", err)).WriteTo(w)
  }
}

func (handler *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  defer errorHandling(w, r)
  var response Response
  switch r.Method {
    case "POST": response = handler.Self.Post((*JsonRequest)(r))
    case "GET": response = handler.Self.Get(r)
    case "PUT": response = handler.Self.Put((*JsonRequest)(r))
    case "DELETE": response = handler.Self.Delete(r)
    default: response = NotFound{}
  }
  switch err := response.(type) {
    case ServerError:
      c := appengine.NewContext(r)
      c.Errorf(string(err))
      response.WriteTo(w)
    default: response.WriteTo(w)
  }
}
func (handler *BaseHandler) Post(r *JsonRequest) Response {
  return NotFound{}
}
func (handler *BaseHandler) Get(r *http.Request) Response {
  return NotFound{}
}
func (handler *BaseHandler) Put(r *JsonRequest) Response {
  return NotFound{}
}
func (handler *BaseHandler) Delete(r *http.Request) Response {
  return NotFound{}
}
