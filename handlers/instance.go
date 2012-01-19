package handlers

import (
  "fmt"
  "http"
  "lib"
  "appengine"
  "rackspace"
)
func init() {
  i := new(instanceHandler)
  i.Self = i
  http.Handle("/instance", i)
}

type instanceHandler struct { lib.BaseHandler }

func (handler *instanceHandler) Post(w lib.JsonResponse, r *lib.JsonRequest) {
  proxy := rackspace.NewProxy(appengine.NewContext((*http.Request)(r)))
  fmt.Print(proxy.ListImages())
}
