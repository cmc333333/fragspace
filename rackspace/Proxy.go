package rackspace

import (
  "fmt"
  "bytes"
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

func (proxy *Proxy) post(addr string, input interface{}, output interface{}) os.Error {
  j, err := json.Marshal(input)
  if err != nil {
    panic(lib.ServerError{"problem marshalling json: " + err.String()})
  }
  req, err := http.NewRequest("POST", proxy.auth.Url + addr, bytes.NewBuffer(j))
  if err != nil {
    panic(lib.ServerError{"could not create request: " + err.String()})
  }
  req.Header.Add("Content-Type", "application/json")
  fmt.Print(string(j))
  return proxy.execute(req, output)
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
  fmt.Print(string(body))
  return json.Unmarshal(body, toObj)
}
