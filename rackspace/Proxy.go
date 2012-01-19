package rackspace

import (
  "io/ioutil"
  "http"
  "json"
  "os"
  "appengine"
  "appengine/urlfetch"
  "lib"
)

type Proxy struct {
  auth Auth
  client http.Client
}

func NewProxy(context appengine.Context) *Proxy {
  return &Proxy{*retrieveAuth(context), *urlfetch.Client(context)}
}

func (proxy *Proxy) execute(req *http.Request, toObj interface{}) os.Error {
  req.Header.Add("x-Auth-Token", proxy.auth.Token)
  res, err := proxy.client.Do(req)
  if err != nil {
    panic(lib.ServerError{"problem contacting hosting provider: " + err.String()})
  }

  body, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    panic(lib.ServerError{"problem reading hosting response: " + err.String()})
  }
  return json.Unmarshal(body, toObj)
}
