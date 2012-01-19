package rackspace

import (
  "lib"
  "http"
)

type FlavorsGet struct {
  Flavors []Flavor
}
type Flavor struct {
  Id int
  Name string
  Ram int
  Disk int
}

func (proxy *Proxy) ListFlavors() *FlavorsGet {
  req, err := http.NewRequest("GET", proxy.auth.Url + "/flavors/detail", nil)
  if err != nil {
    panic(lib.ServerError{"could not create request: " + err.String()})
  }
  flavors := new(FlavorsGet)
  proxy.execute(req, flavors)
  return flavors
}
