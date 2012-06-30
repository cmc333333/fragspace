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
  images := proxy.ListImages()
  flavors := proxy.ListFlavors()

  image, flavor := -1, -1
  for _, i := range images.Images {
    if i.Name == "Debian 6 (Squeeze)" {
      image = i.Id
    }
  }
  for _, f := range flavors.Flavors {
    if f.Ram == 256 {
      flavor = f.Id
    }
  }

  fmt.Print(proxy.NewServer(image, flavor))
}
