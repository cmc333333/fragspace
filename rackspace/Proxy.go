package rackspace

import (
  "appengine"
)

type Proxy struct {
  Auth *Auth
}

func NewProxy(context appengine.Context) *Proxy {
  return &Proxy{retrieveAuth(context)}
}
