package rackspace

import (
  "lib"
  "http"
)

type ImagesGet struct {
  Images []Image
}
type Image struct {
  Id int
  Name string
  Status string
}

func (proxy *Proxy) ListImages() *ImagesGet {
  req, err := http.NewRequest("GET", proxy.auth.Url + "/images/detail", nil)
  if err != nil {
    panic(lib.ServerError{"could not create request: " + err.String()})
  }
  images := new(ImagesGet)
  proxy.execute(req, images)
  return images
}
